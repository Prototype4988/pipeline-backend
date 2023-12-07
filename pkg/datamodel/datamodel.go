package datamodel

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"google.golang.org/protobuf/types/known/structpb"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	taskPB "github.com/instill-ai/protogen-go/common/task/v1alpha"
	pipelinePB "github.com/instill-ai/protogen-go/vdp/pipeline/v1beta"
)

// BaseStatic contains common columns for all tables with static UUID as primary key
type BaseStatic struct {
	UID        uuid.UUID      `gorm:"type:uuid;primary_key;<-:create"` // allow read and create
	CreateTime time.Time      `gorm:"autoCreateTime:nano"`
	UpdateTime time.Time      `gorm:"autoUpdateTime:nano"`
	DeleteTime gorm.DeletedAt `sql:"index"`
}

// BaseDynamic contains common columns for all tables with dynamic UUID as primary key generated when creating
type BaseDynamic struct {
	UID        uuid.UUID      `gorm:"type:uuid;primary_key;<-:create"` // allow read and create
	CreateTime time.Time      `gorm:"autoCreateTime:nano"`
	UpdateTime time.Time      `gorm:"autoUpdateTime:nano"`
	DeleteTime gorm.DeletedAt `sql:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *BaseDynamic) BeforeCreate(db *gorm.DB) error {
	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	db.Statement.SetColumn("UID", uuid)
	return nil
}

// Pipeline is the data model of the pipeline table
type Pipeline struct {
	BaseDynamic
	ID                string
	Owner             string
	Description       sql.NullString
	Recipe            *Recipe `gorm:"type:jsonb"`
	DefaultReleaseUID uuid.UUID
	Permission        *Permission `gorm:"type:jsonb"`
	ShareCode         string
	Metadata          datatypes.JSON `gorm:"type:jsonb"`
	Readme            string
}

// PipelineRelease is the data model of the pipeline release table
type PipelineRelease struct {
	BaseDynamic
	ID          string
	PipelineUID uuid.UUID
	Description sql.NullString
	Recipe      *Recipe        `gorm:"type:jsonb"`
	Metadata    datatypes.JSON `gorm:"type:jsonb"`
	Readme      string
}

// Recipe is the data model of the pipeline recipe
type Recipe struct {
	Version    string       `json:"version,omitempty"`
	Components []*Component `json:"components,omitempty"`
}

type Component struct {
	Id             string           `json:"id"`
	DefinitionName string           `json:"definition_name"`
	ResourceName   string           `json:"resource_name"`
	Configuration  *structpb.Struct `json:"configuration"`
}

type Permission struct {
	Users     map[string]*PermissionUser `json:"users,omitempty"`
	ShareCode *PermissionCode            `json:"share_code,omitempty"`
}

// Permission
type PermissionUser struct {
	Enabled bool   `json:"enabled,omitempty"`
	Role    string `json:"role,omitempty"`
}

type PermissionCode struct {
	User    string `json:"user"`
	Code    string `json:"code"`
	Enabled bool   `json:"enabled,omitempty"`
	Role    string `json:"role,omitempty"`
}

// PipelineRole is an alias type for Protobuf enum
type PipelineRole pipelinePB.Role

// Scan function for custom GORM type Recipe
func (r *Recipe) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal Recipe value:", value))
	}

	if err := json.Unmarshal(bytes, &r); err != nil {
		return err
	}

	return nil
}

// Value function for custom GORM type Recipe
func (r *Recipe) Value() (driver.Value, error) {
	valueString, err := json.Marshal(r)
	return string(valueString), err
}

// Scan function for custom GORM type Recipe
func (p *Permission) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal value:", value))
	}

	if err := json.Unmarshal(bytes, &p); err != nil {
		return err
	}

	return nil
}

// Value function for custom GORM type Recipe
func (p *Permission) Value() (driver.Value, error) {
	valueString, err := json.Marshal(p)
	return string(valueString), err
}

// Scan function for custom GORM type Recipe
func (p *PermissionUser) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal value:", value))
	}

	if err := json.Unmarshal(bytes, &p); err != nil {
		return err
	}

	return nil
}

// Value function for custom GORM type Recipe
func (p *PermissionUser) Value() (driver.Value, error) {
	valueString, err := json.Marshal(p)
	return string(valueString), err
}

// Scan function for custom GORM type Recipe
func (p *PermissionCode) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal value:", value))
	}

	if err := json.Unmarshal(bytes, &p); err != nil {
		return err
	}

	return nil
}

// Value function for custom GORM type Recipe
func (p *PermissionCode) Value() (driver.Value, error) {
	valueString, err := json.Marshal(p)
	return string(valueString), err
}

// Scan function for custom GORM type ReleaseStage
func (p *PipelineRole) Scan(value interface{}) error {
	*p = PipelineRole(pipelinePB.Role_value[value.(string)])
	return nil
}

// Value function for custom GORM type ReleaseStage
func (p PipelineRole) Value() (driver.Value, error) {
	return pipelinePB.Role(p).String(), nil
}

// Connector is the data model of the connector table
type Connector struct {
	BaseDynamic
	ID                     string
	Owner                  string
	ConnectorDefinitionUID uuid.UUID
	Description            sql.NullString
	Tombstone              bool
	Configuration          datatypes.JSON      `gorm:"type:jsonb"`
	ConnectorType          ConnectorType       `sql:"type:CONNECTOR_VALID_CONNECTOR_TYPE"`
	State                  ConnectorState      `sql:"type:CONNECTOR_VALID_CONNECTOR_STATE"`
	Visibility             ConnectorVisibility `sql:"type:CONNECTOR_VALID_CONNECTOR_VISIBILITY"`
}

func (Connector) TableName() string {
	return "connector"
}

// ConnectorType is an alias type for Protobuf enum ConnectorType
type ConnectorVisibility pipelinePB.Connector_Visibility

// ConnectorType is an alias type for Protobuf enum ConnectorType
type ConnectorType pipelinePB.ConnectorType

// ConnectorType is an alias type for Protobuf enum ConnectorType
type Task taskPB.Task

// Scan function for custom GORM type ConnectorType
func (c *ConnectorType) Scan(value interface{}) error {
	*c = ConnectorType(pipelinePB.ConnectorType_value[value.(string)])
	return nil
}

// Value function for custom GORM type ConnectorType
func (c ConnectorType) Value() (driver.Value, error) {
	return pipelinePB.ConnectorType(c).String(), nil
}

// ConnectorState is an alias type for Protobuf enum ConnectorState
type ConnectorState pipelinePB.Connector_State

// Scan function for custom GORM type ConnectorState
func (r *ConnectorState) Scan(value interface{}) error {
	*r = ConnectorState(pipelinePB.Connector_State_value[value.(string)])
	return nil
}

// Value function for custom GORM type ConnectorState
func (r ConnectorState) Value() (driver.Value, error) {
	return pipelinePB.Connector_State(r).String(), nil
}

// Scan function for custom GORM type ReleaseStage
func (r *ConnectorVisibility) Scan(value interface{}) error {
	*r = ConnectorVisibility(pipelinePB.Connector_Visibility_value[value.(string)])
	return nil
}

// Value function for custom GORM type ReleaseStage
func (r ConnectorVisibility) Value() (driver.Value, error) {
	return pipelinePB.Connector_Visibility(r).String(), nil
}

// Scan function for custom GORM type Task
func (r *Task) Scan(value interface{}) error {
	*r = Task(taskPB.Task_value[value.(string)])
	return nil
}

// Value function for custom GORM type Task
func (r Task) Value() (driver.Value, error) {
	return taskPB.Task(r).String(), nil
}
