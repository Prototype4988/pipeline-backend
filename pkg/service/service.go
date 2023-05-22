package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/gofrs/uuid"
	"github.com/gogo/status"
	"github.com/oklog/ulid/v2"
	"go.einride.tech/aip/filtering"
	"google.golang.org/grpc/codes"

	"github.com/instill-ai/pipeline-backend/internal/resource"
	"github.com/instill-ai/pipeline-backend/pkg/datamodel"
	"github.com/instill-ai/pipeline-backend/pkg/logger"
	"github.com/instill-ai/pipeline-backend/pkg/repository"

	connectorPB "github.com/instill-ai/protogen-go/vdp/connector/v1alpha"
	controllerPB "github.com/instill-ai/protogen-go/vdp/controller/v1alpha"
	mgmtPB "github.com/instill-ai/protogen-go/vdp/mgmt/v1alpha"
	modelPB "github.com/instill-ai/protogen-go/vdp/model/v1alpha"
	pipelinePB "github.com/instill-ai/protogen-go/vdp/pipeline/v1alpha"
)

type TextToImageInput struct {
	Prompt   string
	Steps    int64
	CfgScale float32
	Seed     int64
	Samples  int64
}

type TextGenerationInput struct {
	Prompt        string
	OutputLen     int64
	BadWordsList  string
	StopWordsList string
	TopK          int64
	Seed          int64
}

type ImageInput struct {
	Content     []byte
	FileNames   []string
	FileLengths []uint64
}

// Service interface
type Service interface {
	GetMgmtPrivateServiceClient() mgmtPB.MgmtPrivateServiceClient
	GetRedisClient() *redis.Client

	CreatePipeline(owner *mgmtPB.User, pipeline *datamodel.Pipeline) (*datamodel.Pipeline, error)
	ListPipelines(owner *mgmtPB.User, pageSize int64, pageToken string, isBasicView bool, filter filtering.Filter) ([]datamodel.Pipeline, int64, string, error)
	GetPipelineByID(id string, owner *mgmtPB.User, isBasicView bool) (*datamodel.Pipeline, error)
	GetPipelineByUID(uid uuid.UUID, owner *mgmtPB.User, isBasicView bool) (*datamodel.Pipeline, error)
	UpdatePipeline(id string, owner *mgmtPB.User, updatedPipeline *datamodel.Pipeline) (*datamodel.Pipeline, error)
	DeletePipeline(id string, owner *mgmtPB.User) error
	UpdatePipelineState(id string, owner *mgmtPB.User, state datamodel.PipelineState) (*datamodel.Pipeline, error)
	UpdatePipelineID(id string, owner *mgmtPB.User, newID string) (*datamodel.Pipeline, error)
	TriggerPipeline(req *pipelinePB.TriggerPipelineRequest, owner *mgmtPB.User, pipeline *datamodel.Pipeline) (*pipelinePB.TriggerPipelineResponse, error)
	TriggerPipelineBinaryFileUpload(owner *mgmtPB.User, pipeline *datamodel.Pipeline, task modelPB.Model_Task, input interface{}) (*pipelinePB.TriggerPipelineBinaryFileUploadResponse, error)
	GetModelByName(owner *mgmtPB.User, modelName string) (*modelPB.Model, error)

	ListPipelinesAdmin(pageSize int64, pageToken string, isBasicView bool, filter filtering.Filter) ([]datamodel.Pipeline, int64, string, error)
	GetPipelineByUIDAdmin(uid uuid.UUID, isBasicView bool) (*datamodel.Pipeline, error)
	// Controller APIs
	GetResourceState(uid uuid.UUID) (*pipelinePB.Pipeline_State, error)
	UpdateResourceState(uid uuid.UUID, state pipelinePB.Pipeline_State, progress *int32) error
	DeleteResourceState(uid uuid.UUID) error
}

type service struct {
	repository                    repository.Repository
	mgmtPrivateServiceClient      mgmtPB.MgmtPrivateServiceClient
	connectorPublicServiceClient  connectorPB.ConnectorPublicServiceClient
	connectorPrivateServiceClient connectorPB.ConnectorPrivateServiceClient
	modelPublicServiceClient      modelPB.ModelPublicServiceClient
	modelPrivateServiceClient     modelPB.ModelPrivateServiceClient
	controllerClient              controllerPB.ControllerPrivateServiceClient
	redisClient                   *redis.Client
}

// NewService initiates a service instance
func NewService(r repository.Repository,
	u mgmtPB.MgmtPrivateServiceClient,
	c connectorPB.ConnectorPublicServiceClient,
	cPrivate connectorPB.ConnectorPrivateServiceClient,
	m modelPB.ModelPublicServiceClient,
	mPrivate modelPB.ModelPrivateServiceClient,
	ct controllerPB.ControllerPrivateServiceClient,
	rc *redis.Client,
) Service {
	return &service{
		repository:                    r,
		mgmtPrivateServiceClient:      u,
		connectorPublicServiceClient:  c,
		connectorPrivateServiceClient: cPrivate,
		modelPublicServiceClient:      m,
		modelPrivateServiceClient:     mPrivate,
		controllerClient:              ct,
		redisClient:                   rc,
	}
}

// GetMgmtPrivateServiceClient returns the management private service client
func (h *service) GetMgmtPrivateServiceClient() mgmtPB.MgmtPrivateServiceClient {
	return h.mgmtPrivateServiceClient
}

// GetRedisClient returns the redis client
func (h *service) GetRedisClient() *redis.Client {
	return h.redisClient
}

func (s *service) CreatePipeline(owner *mgmtPB.User, dbPipeline *datamodel.Pipeline) (*datamodel.Pipeline, error) {

	mode, err := s.checkRecipe(owner, dbPipeline.Recipe)
	if err != nil {
		return nil, err
	}

	dbPipeline.Mode = mode

	// User desires to be active
	dbPipeline.State = datamodel.PipelineState(pipelinePB.Pipeline_STATE_ACTIVE)

	ownerPermalink := GenOwnerPermalink(owner)
	dbPipeline.Owner = ownerPermalink

	recipePermalink, err := s.recipeNameToPermalink(owner, dbPipeline.Recipe)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	dbPipeline.Recipe = recipePermalink

	resourceState, err := s.checkState(dbPipeline.Recipe)
	if err != nil {
		return nil, err
	}

	if err := s.repository.CreatePipeline(dbPipeline); err != nil {
		return nil, err
	}

	dbCreatedPipeline, err := s.repository.GetPipelineByID(dbPipeline.ID, ownerPermalink, false)
	if err != nil {
		return nil, err
	}

	rErr := s.includeResourceDetailInRecipe(dbCreatedPipeline.Recipe)
	if rErr != nil {
		return nil, rErr
	}
	createdCecipeRscName, err := s.recipePermalinkToName(dbCreatedPipeline.Recipe)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// Add resource entry to controller to start checking components' state
	if err := s.UpdateResourceState(dbCreatedPipeline.UID, pipelinePB.Pipeline_State(resourceState), nil); err != nil {
		return nil, err
	}

	dbCreatedPipeline.Recipe = createdCecipeRscName

	return dbCreatedPipeline, nil
}

func (s *service) ListPipelines(owner *mgmtPB.User, pageSize int64, pageToken string, isBasicView bool, filter filtering.Filter) ([]datamodel.Pipeline, int64, string, error) {

	ownerPermalink := GenOwnerPermalink(owner)
	dbPipelines, ps, pt, err := s.repository.ListPipelines(ownerPermalink, pageSize, pageToken, isBasicView, filter)
	if err != nil {
		return nil, 0, "", err
	}

	if !isBasicView {
		for idx := range dbPipelines {
			err := s.includeResourceDetailInRecipe(dbPipelines[idx].Recipe)
			if err != nil {
				return nil, 0, "", err
			}
			recipeRscName, err := s.recipePermalinkToName(dbPipelines[idx].Recipe)
			if err != nil {
				return nil, 0, "", status.Errorf(codes.Internal, err.Error())
			}
			dbPipelines[idx].Recipe = recipeRscName
		}
	}

	return dbPipelines, ps, pt, nil
}

func (s *service) ListPipelinesAdmin(pageSize int64, pageToken string, isBasicView bool, filter filtering.Filter) ([]datamodel.Pipeline, int64, string, error) {

	dbPipelines, ps, pt, err := s.repository.ListPipelinesAdmin(pageSize, pageToken, isBasicView, filter)
	if err != nil {
		return nil, 0, "", err
	}
	if !isBasicView {
		for idx := range dbPipelines {
			err := s.includeResourceDetailInRecipe(dbPipelines[idx].Recipe)
			if err != nil {
				return nil, 0, "", err
			}
		}
	}

	return dbPipelines, ps, pt, nil
}

func (s *service) GetPipelineByID(id string, owner *mgmtPB.User, isBasicView bool) (*datamodel.Pipeline, error) {

	ownerPermalink := GenOwnerPermalink(owner)

	dbPipeline, err := s.repository.GetPipelineByID(id, ownerPermalink, isBasicView)
	if err != nil {
		return nil, err
	}

	if !isBasicView {
		err := s.includeResourceDetailInRecipe(dbPipeline.Recipe)
		if err != nil {
			return nil, err
		}
		recipeRscName, err := s.recipePermalinkToName(dbPipeline.Recipe)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		dbPipeline.Recipe = recipeRscName
	}

	return dbPipeline, nil
}

func (s *service) GetPipelineByUID(uid uuid.UUID, owner *mgmtPB.User, isBasicView bool) (*datamodel.Pipeline, error) {

	ownerPermalink := GenOwnerPermalink(owner)

	dbPipeline, err := s.repository.GetPipelineByUID(uid, ownerPermalink, isBasicView)
	if err != nil {
		return nil, err
	}

	if !isBasicView {
		err := s.includeResourceDetailInRecipe(dbPipeline.Recipe)
		if err != nil {
			return nil, err
		}
		recipeRscName, err := s.recipePermalinkToName(dbPipeline.Recipe)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		dbPipeline.Recipe = recipeRscName
	}

	return dbPipeline, nil
}

func (s *service) GetPipelineByUIDAdmin(uid uuid.UUID, isBasicView bool) (*datamodel.Pipeline, error) {

	dbPipeline, err := s.repository.GetPipelineByUIDAdmin(uid, isBasicView)
	if err != nil {
		return nil, err
	}
	if !isBasicView {
		err := s.includeResourceDetailInRecipe(dbPipeline.Recipe)
		if err != nil {
			return nil, err
		}
	}

	return dbPipeline, nil
}

func (s *service) UpdatePipeline(id string, owner *mgmtPB.User, toUpdPipeline *datamodel.Pipeline) (*datamodel.Pipeline, error) {

	resourceState := toUpdPipeline.State

	if toUpdPipeline.Recipe != nil {
		mode, err := s.checkRecipe(owner, toUpdPipeline.Recipe)
		if err != nil {
			return nil, err
		}

		toUpdPipeline.Mode = mode

		// User desires to be active
		toUpdPipeline.State = datamodel.PipelineState(pipelinePB.Pipeline_STATE_ACTIVE)

		recipePermalink, err := s.recipeNameToPermalink(owner, toUpdPipeline.Recipe)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}

		toUpdPipeline.Recipe = recipePermalink

		resourceState, err = s.checkState(toUpdPipeline.Recipe)
		if err != nil {
			return nil, err
		}
	}

	ownerPermalink := GenOwnerPermalink(owner)

	toUpdPipeline.Owner = ownerPermalink

	// Validation: Pipeline existence
	if existingPipeline, _ := s.repository.GetPipelineByID(id, ownerPermalink, true); existingPipeline == nil {
		return nil, status.Errorf(codes.NotFound, "Pipeline id %s is not found", id)
	}

	if err := s.repository.UpdatePipeline(id, ownerPermalink, toUpdPipeline); err != nil {
		return nil, err
	}

	dbPipeline, err := s.repository.GetPipelineByID(toUpdPipeline.ID, ownerPermalink, false)
	if err != nil {
		return nil, err
	}

	// Update resource entry in controller to start checking components' state
	if err := s.UpdateResourceState(dbPipeline.UID, pipelinePB.Pipeline_State(resourceState), nil); err != nil {
		return nil, err
	}

	rErr := s.includeResourceDetailInRecipe(dbPipeline.Recipe)
	if rErr != nil {
		return nil, rErr
	}

	recipeRscName, err := s.recipePermalinkToName(dbPipeline.Recipe)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	dbPipeline.Recipe = recipeRscName

	return dbPipeline, nil
}

func (s *service) DeletePipeline(id string, owner *mgmtPB.User) error {
	ownerPermalink := GenOwnerPermalink(owner)

	dbPipeline, err := s.repository.GetPipelineByID(id, ownerPermalink, false)
	if err != nil {
		return err
	}

	if err := s.DeleteResourceState(dbPipeline.UID); err != nil {
		return err
	}

	return s.repository.DeletePipeline(id, ownerPermalink)
}

func (s *service) UpdatePipelineState(id string, owner *mgmtPB.User, state datamodel.PipelineState) (*datamodel.Pipeline, error) {

	if state == datamodel.PipelineState(pipelinePB.Pipeline_STATE_UNSPECIFIED) {
		return nil, status.Errorf(codes.InvalidArgument, "State update with unspecified is not allowed")
	}

	ownerPermalink := GenOwnerPermalink(owner)

	dbPipeline, err := s.repository.GetPipelineByID(id, ownerPermalink, false)
	if err != nil {
		return nil, err
	}

	// user desires to be active or inactive, state stay the same
	// but update etcd storage with checkState
	var resourceState datamodel.PipelineState
	if state == datamodel.PipelineState(pipelinePB.Pipeline_STATE_ACTIVE) {
		resourceState, err = s.checkState(dbPipeline.Recipe)
		if err != nil {
			return nil, err
		}
	} else {
		resourceState = datamodel.PipelineState(pipelinePB.Pipeline_STATE_INACTIVE)
	}

	recipeRscName, err := s.recipePermalinkToName(dbPipeline.Recipe)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	mode, recipeErr := s.checkRecipe(owner, recipeRscName)
	if recipeErr != nil {
		return nil, err
	}

	if mode == datamodel.PipelineMode(pipelinePB.Pipeline_MODE_SYNC) && state == datamodel.PipelineState(pipelinePB.Pipeline_STATE_INACTIVE) {
		return nil, status.Errorf(codes.InvalidArgument, "Pipeline %s is in the SYNC mode, which is always active", dbPipeline.ID)
	}

	if err := s.repository.UpdatePipelineState(id, ownerPermalink, state); err != nil {
		return nil, err
	}

	dbPipeline, err = s.repository.GetPipelineByID(id, ownerPermalink, false)
	if err != nil {
		return nil, err
	}

	// Update resource entry in controller to start checking components' state
	if err := s.UpdateResourceState(dbPipeline.UID, pipelinePB.Pipeline_State(resourceState), nil); err != nil {
		return nil, err
	}

	dbPipeline.Recipe = recipeRscName

	return dbPipeline, nil
}

func (s *service) UpdatePipelineID(id string, owner *mgmtPB.User, newID string) (*datamodel.Pipeline, error) {

	ownerPermalink := GenOwnerPermalink(owner)

	// Validation: Pipeline existence
	existingPipeline, _ := s.repository.GetPipelineByID(id, ownerPermalink, true)
	if existingPipeline == nil {
		return nil, status.Errorf(codes.NotFound, "Pipeline id %s is not found", id)
	}

	// if err := s.DeleteResourceState(existingPipeline.UID.String()); err != nil {
	// 	return nil, err
	// }

	if err := s.repository.UpdatePipelineID(id, ownerPermalink, newID); err != nil {
		return nil, err
	}

	dbPipeline, err := s.repository.GetPipelineByID(newID, ownerPermalink, false)
	if err != nil {
		return nil, err
	}

	// if err := s.UpdateResourceState(dbPipeline.UID.String(), pipelinePB.Pipeline_State(existingPipeline.State), nil); err != nil {
	// 	return nil, err
	// }

	recipeRscName, err := s.recipePermalinkToName(dbPipeline.Recipe)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	dbPipeline.Recipe = recipeRscName

	return dbPipeline, nil
}

func (s *service) TriggerPipeline(req *pipelinePB.TriggerPipelineRequest, owner *mgmtPB.User, dbPipeline *datamodel.Pipeline) (*pipelinePB.TriggerPipelineResponse, error) {

	logger, _ := logger.GetZapLogger()

	if dbPipeline.State != datamodel.PipelineState(pipelinePB.Pipeline_STATE_ACTIVE) {
		return nil, status.Error(codes.FailedPrecondition, fmt.Sprintf("The pipeline %s is not active", dbPipeline.ID))
	}

	var dataMappingIndices []string
	var taskInputs []*modelPB.TaskInput
	for _, taskInput := range req.TaskInputs {
		taskInputs = append(taskInputs, &modelPB.TaskInput{
			Input: taskInput.GetInput(),
		})
		dataMappingIndices = append(dataMappingIndices, ulid.Make().String())
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	models := GetModelsFromRecipe(dbPipeline.Recipe)
	dests := GetDestinationsFromRecipe(dbPipeline.Recipe)
	var modelOutputs []*pipelinePB.ModelOutput
	errors := make(chan error)
	go func() {
		defer wg.Done()

		for idx, model := range models {

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
			defer cancel()
			resp, err := s.modelPublicServiceClient.TriggerModel(InjectOwnerToContext(ctx, owner), &modelPB.TriggerModelRequest{
				Name:       model,
				TaskInputs: taskInputs,
			})
			if err != nil {
				errors <- err
				logger.Error(fmt.Sprintf("[model-backend] Error %s at %dth model %s: %v", "TriggerModel", idx, model, err.Error()))
				return
			}

			taskOutputs := cvtModelTaskOutputToPipelineTaskOutput(resp.TaskOutputs)
			for idx, taskOutput := range taskOutputs {
				taskOutput.Index = dataMappingIndices[idx]
			}

			modelOutputs = append(modelOutputs, &pipelinePB.ModelOutput{
				Model:       model,
				Task:        resp.Task,
				TaskOutputs: taskOutputs,
			})

			// Increment trigger image numbers
			uid, err := resource.GetPermalinkUID(dbPipeline.Owner)
			if err != nil {
				errors <- err
				logger.Error(err.Error())
				return
			}
			if strings.HasPrefix(dbPipeline.Owner, "users/") {
				s.redisClient.IncrBy(context.Background(), fmt.Sprintf("user:%s:trigger.num", uid), int64(len(taskInputs)))
			} else if strings.HasPrefix(dbPipeline.Owner, "orgs/") {
				s.redisClient.IncrBy(context.Background(), fmt.Sprintf("org:%s:trigger.num", uid), int64(len(taskInputs)))
			}
		}
		errors <- nil
	}()

	switch {
	// If this is a SYNC trigger (i.e., HTTP, gRPC source and destination connectors), return right away
	case dbPipeline.Mode == datamodel.PipelineMode(pipelinePB.Pipeline_MODE_SYNC):
		go func() {
			wg.Wait()
			close(errors)
		}()
		for err := range errors {
			if err != nil {
				switch {
				case strings.Contains(err.Error(), "code = DeadlineExceeded"):
					return nil, status.Errorf(codes.DeadlineExceeded, "trigger model got timeout error")
				default:
					return nil, status.Errorf(codes.Internal, fmt.Sprintf("trigger model got error %v", err.Error()))
				}
			}
		}
		return &pipelinePB.TriggerPipelineResponse{
			DataMappingIndices: dataMappingIndices,
			ModelOutputs:       modelOutputs,
		}, nil
	// If this is a async trigger, write to the destination connector
	case dbPipeline.Mode == datamodel.PipelineMode(pipelinePB.Pipeline_MODE_ASYNC):
		go func() {
			go func() {
				wg.Wait()
				close(errors)
			}()
			for err := range errors {
				if err != nil {
					logger.Error(fmt.Sprintf("[model-backend] Error trigger model got error %v", err.Error()))
					return
				}
			}

			for idx, destRecName := range dests {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				_, err := s.connectorPublicServiceClient.WriteDestinationConnector(InjectOwnerToContext(ctx, owner), &connectorPB.WriteDestinationConnectorRequest{
					Name:                destRecName,
					SyncMode:            connectorPB.SupportedSyncModes_SUPPORTED_SYNC_MODES_FULL_REFRESH,
					DestinationSyncMode: connectorPB.SupportedDestinationSyncModes_SUPPORTED_DESTINATION_SYNC_MODES_APPEND,
					Pipeline:            fmt.Sprintf("pipelines/%s", dbPipeline.ID),
					DataMappingIndices:  dataMappingIndices,
					ModelOutputs:        modelOutputs,
					Recipe: func() *pipelinePB.Recipe {
						if dbPipeline.Recipe != nil {
							if err := s.excludeResourceDetailFromRecipe(dbPipeline.Recipe); err != nil {
								logger.Error(err.Error())
							}
							b, err := json.Marshal(dbPipeline.Recipe)
							if err != nil {
								logger.Error(err.Error())
							}
							pbRecipe := pipelinePB.Recipe{}
							err = json.Unmarshal(b, &pbRecipe)
							if err != nil {
								logger.Error(err.Error())
							}
							return &pbRecipe
						}
						return nil
					}(),
				})
				if err != nil {
					logger.Error(fmt.Sprintf("[connector-backend] Error %s at %dth destination %s: %v", "WriteDestinationConnector", idx, destRecName, err.Error()))
				}
			}
		}()
		return &pipelinePB.TriggerPipelineResponse{
			DataMappingIndices: dataMappingIndices,
			ModelOutputs:       nil,
		}, nil
	}

	return nil, status.Errorf(codes.Internal, "something went very wrong - unable to trigger the pipeline")

}

func (s *service) triggerImageTask(owner *mgmtPB.User, dbPipeline *datamodel.Pipeline, task modelPB.Model_Task, input interface{}, dataMappingIndices []string) ([]*pipelinePB.ModelOutput, error) {
	imageInput := input.(*ImageInput)
	var modelOutputs []*pipelinePB.ModelOutput

	models := GetModelsFromRecipe(dbPipeline.Recipe)

	for idx, model := range models {

		// TODO: async call model-backend
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		stream, err := s.modelPublicServiceClient.TriggerModelBinaryFileUpload(InjectOwnerToContext(ctx, owner))
		defer func() {
			_ = stream.CloseSend()
		}()

		if err != nil {
			return modelOutputs, status.Errorf(codes.Internal, "[model-backend] Error %s at %dth model %s: cannot init stream: %v", "TriggerModelBinaryFileUpload", idx, model, err.Error())
		}
		var triggerRequest modelPB.TriggerModelBinaryFileUploadRequest
		switch task {
		case modelPB.Model_TASK_CLASSIFICATION:
			triggerRequest = modelPB.TriggerModelBinaryFileUploadRequest{
				Name: model,
				TaskInput: &modelPB.TaskInputStream{
					Input: &modelPB.TaskInputStream_Classification{
						Classification: &modelPB.ClassificationInputStream{
							FileLengths: imageInput.FileLengths,
						},
					},
				},
			}
		case modelPB.Model_TASK_DETECTION:
			triggerRequest = modelPB.TriggerModelBinaryFileUploadRequest{
				Name: model,
				TaskInput: &modelPB.TaskInputStream{
					Input: &modelPB.TaskInputStream_Detection{
						Detection: &modelPB.DetectionInputStream{
							FileLengths: imageInput.FileLengths,
						},
					},
				},
			}
		case modelPB.Model_TASK_KEYPOINT:
			triggerRequest = modelPB.TriggerModelBinaryFileUploadRequest{
				Name: model,
				TaskInput: &modelPB.TaskInputStream{
					Input: &modelPB.TaskInputStream_Keypoint{
						Keypoint: &modelPB.KeypointInputStream{
							FileLengths: imageInput.FileLengths,
						},
					},
				},
			}
		case modelPB.Model_TASK_OCR:
			triggerRequest = modelPB.TriggerModelBinaryFileUploadRequest{
				Name: model,
				TaskInput: &modelPB.TaskInputStream{
					Input: &modelPB.TaskInputStream_Ocr{
						Ocr: &modelPB.OcrInputStream{
							FileLengths: imageInput.FileLengths,
						},
					},
				},
			}

		case modelPB.Model_TASK_INSTANCE_SEGMENTATION:
			triggerRequest = modelPB.TriggerModelBinaryFileUploadRequest{
				Name: model,
				TaskInput: &modelPB.TaskInputStream{
					Input: &modelPB.TaskInputStream_InstanceSegmentation{
						InstanceSegmentation: &modelPB.InstanceSegmentationInputStream{
							FileLengths: imageInput.FileLengths,
						},
					},
				},
			}
		case modelPB.Model_TASK_SEMANTIC_SEGMENTATION:
			triggerRequest = modelPB.TriggerModelBinaryFileUploadRequest{
				Name: model,
				TaskInput: &modelPB.TaskInputStream{
					Input: &modelPB.TaskInputStream_SemanticSegmentation{
						SemanticSegmentation: &modelPB.SemanticSegmentationInputStream{
							FileLengths: imageInput.FileLengths,
						},
					},
				},
			}
		}
		if err := stream.Send(&triggerRequest); err != nil {
			return modelOutputs, status.Errorf(codes.Internal, "[model-backend] Error %s at %dth model %s: cannot send data info to server: %v", "TriggerModelBinaryFileUploadRequest", idx, model, err.Error())
		}
		fb := bytes.Buffer{}
		fb.Write(imageInput.Content)
		buf := make([]byte, 64*1024)
		for {
			n, err := fb.Read(buf)
			if err == io.EOF {
				break
			} else if err != nil {
				return modelOutputs, err
			}

			var triggerRequest modelPB.TriggerModelBinaryFileUploadRequest
			switch task {
			case modelPB.Model_TASK_CLASSIFICATION:
				triggerRequest = modelPB.TriggerModelBinaryFileUploadRequest{
					TaskInput: &modelPB.TaskInputStream{
						Input: &modelPB.TaskInputStream_Classification{
							Classification: &modelPB.ClassificationInputStream{
								Content: buf[:n],
							},
						},
					},
				}
			case modelPB.Model_TASK_DETECTION:
				triggerRequest = modelPB.TriggerModelBinaryFileUploadRequest{
					TaskInput: &modelPB.TaskInputStream{
						Input: &modelPB.TaskInputStream_Detection{
							Detection: &modelPB.DetectionInputStream{
								Content: buf[:n],
							},
						},
					},
				}
			case modelPB.Model_TASK_KEYPOINT:
				triggerRequest = modelPB.TriggerModelBinaryFileUploadRequest{
					TaskInput: &modelPB.TaskInputStream{
						Input: &modelPB.TaskInputStream_Keypoint{
							Keypoint: &modelPB.KeypointInputStream{
								Content: buf[:n],
							},
						},
					},
				}
			case modelPB.Model_TASK_OCR:
				triggerRequest = modelPB.TriggerModelBinaryFileUploadRequest{
					TaskInput: &modelPB.TaskInputStream{
						Input: &modelPB.TaskInputStream_Ocr{
							Ocr: &modelPB.OcrInputStream{
								Content: buf[:n],
							},
						},
					},
				}

			case modelPB.Model_TASK_INSTANCE_SEGMENTATION:
				triggerRequest = modelPB.TriggerModelBinaryFileUploadRequest{
					TaskInput: &modelPB.TaskInputStream{
						Input: &modelPB.TaskInputStream_InstanceSegmentation{
							InstanceSegmentation: &modelPB.InstanceSegmentationInputStream{
								Content: buf[:n],
							},
						},
					},
				}
			case modelPB.Model_TASK_SEMANTIC_SEGMENTATION:
				triggerRequest = modelPB.TriggerModelBinaryFileUploadRequest{
					TaskInput: &modelPB.TaskInputStream{
						Input: &modelPB.TaskInputStream_SemanticSegmentation{
							SemanticSegmentation: &modelPB.SemanticSegmentationInputStream{
								Content: buf[:n],
							},
						},
					},
				}
			}
			if err := stream.Send(&triggerRequest); err != nil {
				return modelOutputs, status.Errorf(codes.Internal, "[model-backend] Error %s at %dth model %s: cannot send chunk to server: %v", "TriggerModelBinaryFileUploadRequest", idx, model, err.Error())
			}
		}

		resp, err := stream.CloseAndRecv()
		if err != nil {
			return modelOutputs, status.Errorf(codes.Internal, "[model-backend] Error %s at %dth model %s: cannot receive response: %v", "TriggerModelBinaryFileUploadRequest", idx, model, err.Error())
		}

		taskOutputs := cvtModelTaskOutputToPipelineTaskOutput(resp.TaskOutputs)
		for idx, taskOutput := range taskOutputs {
			taskOutput.Index = dataMappingIndices[idx]
		}

		modelOutputs = append(modelOutputs, &pipelinePB.ModelOutput{
			Model:       model,
			Task:        resp.Task,
			TaskOutputs: taskOutputs,
		})

		// Increment trigger image numbers
		uid, err := resource.GetPermalinkUID(dbPipeline.Owner)
		if err != nil {
			return modelOutputs, err
		}
		if strings.HasPrefix(dbPipeline.Owner, "users/") {
			s.redisClient.IncrBy(ctx, fmt.Sprintf("user:%s:trigger.num", uid), int64(len(imageInput.FileLengths)))
		} else if strings.HasPrefix(dbPipeline.Owner, "orgs/") {
			s.redisClient.IncrBy(ctx, fmt.Sprintf("org:%s:trigger.num", uid), int64(len(imageInput.FileLengths)))
		}
	}

	return modelOutputs, nil
}

func (s *service) triggerTextTask(owner *mgmtPB.User, dbPipeline *datamodel.Pipeline, task modelPB.Model_Task, input interface{}, dataMappingIndices []string) ([]*pipelinePB.ModelOutput, error) {

	var modelOutputs []*pipelinePB.ModelOutput
	models := GetModelsFromRecipe(dbPipeline.Recipe)

	for idx, model := range models {

		// TODO: async call model-backend
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		stream, err := s.modelPublicServiceClient.TriggerModelBinaryFileUpload(InjectOwnerToContext(ctx, owner))
		defer func() {
			_ = stream.CloseSend()
		}()

		if err != nil {
			return modelOutputs, status.Errorf(codes.Internal, "[model-backend] Error %s at %dth model %s: cannot init stream: %v", "TriggerModelBinaryFileUpload", idx, model, err.Error())
		}

		var triggerRequest modelPB.TriggerModelBinaryFileUploadRequest
		switch task {
		case modelPB.Model_TASK_TEXT_TO_IMAGE:
			textToImageInput := input.(*TextToImageInput)
			triggerRequest = modelPB.TriggerModelBinaryFileUploadRequest{
				Name: model,
				TaskInput: &modelPB.TaskInputStream{
					Input: &modelPB.TaskInputStream_TextToImage{
						TextToImage: &modelPB.TextToImageInput{
							Prompt:   textToImageInput.Prompt,
							Steps:    &textToImageInput.Steps,
							CfgScale: &textToImageInput.CfgScale,
							Seed:     &textToImageInput.Seed,
							Samples:  &textToImageInput.Samples,
						},
					},
				},
			}
		case modelPB.Model_TASK_TEXT_GENERATION:
			textGenerationInput := input.(*TextGenerationInput)
			triggerRequest = modelPB.TriggerModelBinaryFileUploadRequest{
				Name: model,
				TaskInput: &modelPB.TaskInputStream{
					Input: &modelPB.TaskInputStream_TextGeneration{
						TextGeneration: &modelPB.TextGenerationInput{
							Prompt:        textGenerationInput.Prompt,
							OutputLen:     &textGenerationInput.OutputLen,
							BadWordsList:  &textGenerationInput.BadWordsList,
							StopWordsList: &textGenerationInput.StopWordsList,
							Topk:          &textGenerationInput.TopK,
							Seed:          &textGenerationInput.Seed,
						},
					},
				},
			}
		}

		if err := stream.Send(&triggerRequest); err != nil {
			return modelOutputs, status.Errorf(codes.Internal, "[model-backend] Error %s at %dth model %s: cannot send data info to server: %v", "TriggerModelBinaryFileUploadRequest", idx, model, err.Error())
		}

		resp, err := stream.CloseAndRecv()
		if err != nil {
			return modelOutputs, status.Errorf(codes.Internal, "[model-backend] Error %s at %dth model %s: cannot receive response: %v", "TriggerModelBinaryFileUploadRequest", idx, model, err.Error())
		}

		taskOutputs := cvtModelTaskOutputToPipelineTaskOutput(resp.TaskOutputs)
		for idx, taskOutput := range taskOutputs {
			taskOutput.Index = dataMappingIndices[idx]
		}

		modelOutputs = append(modelOutputs, &pipelinePB.ModelOutput{
			Model:       model,
			Task:        resp.Task,
			TaskOutputs: taskOutputs,
		})

		// Increment trigger image numbers
		uid, err := resource.GetPermalinkUID(dbPipeline.Owner)
		if err != nil {
			return modelOutputs, err
		}
		if strings.HasPrefix(dbPipeline.Owner, "users/") {
			s.redisClient.IncrBy(ctx, fmt.Sprintf("user:%s:trigger.num", uid), 1)
		} else if strings.HasPrefix(dbPipeline.Owner, "orgs/") {
			s.redisClient.IncrBy(ctx, fmt.Sprintf("org:%s:trigger.num", uid), 1)
		}
	}

	return modelOutputs, nil
}

func (s *service) TriggerPipelineBinaryFileUpload(owner *mgmtPB.User, dbPipeline *datamodel.Pipeline, task modelPB.Model_Task, input interface{}) (*pipelinePB.TriggerPipelineBinaryFileUploadResponse, error) {
	if dbPipeline.State != datamodel.PipelineState(pipelinePB.Pipeline_STATE_ACTIVE) {
		return nil, status.Error(codes.FailedPrecondition, fmt.Sprintf("The pipeline %s is not active", dbPipeline.ID))
	}

	batching := 1
	switch task {
	case modelPB.Model_TASK_CLASSIFICATION,
		modelPB.Model_TASK_DETECTION,
		modelPB.Model_TASK_KEYPOINT,
		modelPB.Model_TASK_OCR,
		modelPB.Model_TASK_INSTANCE_SEGMENTATION,
		modelPB.Model_TASK_SEMANTIC_SEGMENTATION:
		inp := input.(*ImageInput)
		batching = len(inp.FileNames)
	case modelPB.Model_TASK_TEXT_TO_IMAGE,
		modelPB.Model_TASK_TEXT_GENERATION:
		batching = 1
	}
	var dataMappingIndices []string
	for i := 0; i < batching; i++ {
		dataMappingIndices = append(dataMappingIndices, ulid.Make().String())
	}

	var modelOutputs []*pipelinePB.ModelOutput
	var err error
	switch task {
	case modelPB.Model_TASK_CLASSIFICATION,
		modelPB.Model_TASK_DETECTION,
		modelPB.Model_TASK_KEYPOINT,
		modelPB.Model_TASK_OCR,
		modelPB.Model_TASK_INSTANCE_SEGMENTATION,
		modelPB.Model_TASK_SEMANTIC_SEGMENTATION:
		modelOutputs, err = s.triggerImageTask(owner, dbPipeline, task, input, dataMappingIndices)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	case modelPB.Model_TASK_TEXT_TO_IMAGE,
		modelPB.Model_TASK_TEXT_GENERATION:
		modelOutputs, err = s.triggerTextTask(owner, dbPipeline, task, input, dataMappingIndices)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	}
	switch {
	// Check if this is a SYNC trigger (i.e., HTTP, gRPC source and destination connectors)
	case dbPipeline.Mode == datamodel.PipelineMode(pipelinePB.Pipeline_MODE_SYNC):
		return &pipelinePB.TriggerPipelineBinaryFileUploadResponse{
			DataMappingIndices: dataMappingIndices,
			ModelOutputs:       modelOutputs,
		}, nil
	case dbPipeline.Mode == datamodel.PipelineMode(pipelinePB.Pipeline_MODE_ASYNC):
		dests := GetDestinationsFromRecipe(dbPipeline.Recipe)

		for idx, destRecName := range dests {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, err := s.connectorPublicServiceClient.WriteDestinationConnector(InjectOwnerToContext(ctx, owner), &connectorPB.WriteDestinationConnectorRequest{
				Name:                destRecName,
				SyncMode:            connectorPB.SupportedSyncModes_SUPPORTED_SYNC_MODES_FULL_REFRESH,
				DestinationSyncMode: connectorPB.SupportedDestinationSyncModes_SUPPORTED_DESTINATION_SYNC_MODES_APPEND,
				Pipeline:            fmt.Sprintf("pipelines/%s", dbPipeline.ID),
				DataMappingIndices:  dataMappingIndices,
				ModelOutputs:        modelOutputs,
				Recipe: func() *pipelinePB.Recipe {
					logger, _ := logger.GetZapLogger()

					if dbPipeline.Recipe != nil {
						if err := s.excludeResourceDetailFromRecipe(dbPipeline.Recipe); err != nil {
							logger.Error(err.Error())
						}
						b, err := json.Marshal(dbPipeline.Recipe)
						if err != nil {
							logger.Error(err.Error())
						}
						pbRecipe := pipelinePB.Recipe{}
						err = json.Unmarshal(b, &pbRecipe)
						if err != nil {
							logger.Error(err.Error())
						}
						return &pbRecipe
					}
					return nil
				}(),
			})
			if err != nil {
				return nil, status.Errorf(codes.Internal, "[connector-backend] Error %s at %dth destination %s: %v", "WriteDestinationConnector", idx, destRecName, err.Error())
			}
		}
		return &pipelinePB.TriggerPipelineBinaryFileUploadResponse{
			DataMappingIndices: dataMappingIndices,
			ModelOutputs:       nil,
		}, nil

	}

	return nil, status.Errorf(codes.Internal, "something went very wrong - unable to trigger the pipeline")

}

func (s *service) GetModelByName(owner *mgmtPB.User, modelName string) (*modelPB.Model, error) {
	modelResq, err := s.modelPublicServiceClient.GetModel(InjectOwnerToContext(context.Background(), owner), &modelPB.GetModelRequest{
		Name: modelName,
	})
	if err != nil {
		return nil, err
	}
	return modelResq.Model, nil
}
