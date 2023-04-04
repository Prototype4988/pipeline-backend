import grpc from 'k6/net/grpc';
import {
  check,
  group
} from "k6";
import {
  randomString
} from "https://jslib.k6.io/k6-utils/1.1.0/index.js";

import * as constant from "./const.js"

const client = new grpc.Client();
client.load(['proto/vdp/pipeline/v1alpha'], 'pipeline_public_service.proto');

export function CheckCreate() {

  group(`Pipelines API: Create a pipeline [with random "jwt-sub" header]`, () => {

    client.connect(constant.pipelineGRPCPublicHost, {
      plaintext: true
    });

    var reqBody = Object.assign({
      id: randomString(63),
      description: randomString(50),
    },
      constant.detSyncHTTPSingleModelRecipe
    )

    // Cannot create a pipeline of a non-exist user
    check(client.invoke('vdp.pipeline.v1alpha.PipelinePublicService/CreatePipeline', {
      pipeline: reqBody
    }, constant.paramsWithJwt), {
      [`[with random "jwt-sub" header] vdp.pipeline.v1alpha.PipelinePublicService/CreatePipeline response StatusUnknown`]: (r) => r.status === grpc.StatusUnknown,
    })

    client.close();
  });
}

export function CheckList() {

  group(`Pipelines API: List pipelines [with random "jwt-sub" header]`, () => {

    client.connect(constant.pipelineGRPCPublicHost, {
      plaintext: true
    });

    // Cannot list pipelines of a non-exist user
    check(client.invoke('vdp.pipeline.v1alpha.PipelinePublicService/ListPipelines', {}, constant.paramsWithJwt), {
      [`[with random "jwt-sub" header] vdp.pipeline.v1alpha.PipelinePublicService/ListPipelines response StatusUnknown`]: (r) => r.status === grpc.StatusUnknown,
    })

    client.close();
  });
}

export function CheckGet() {

  group(`Pipelines API: Get a pipeline [with random "jwt-sub" header]`, () => {

    client.connect(constant.pipelineGRPCPublicHost, {
      plaintext: true
    });

    var reqBody = Object.assign({
      id: randomString(10),
      description: randomString(50),
    },
      constant.detSyncHTTPSingleModelRecipe
    )

    check(client.invoke('vdp.pipeline.v1alpha.PipelinePublicService/CreatePipeline', {
      pipeline: reqBody
    }), {
      [`vdp.pipeline.v1alpha.PipelinePublicService/CreatePipeline response StatusOK`]: (r) => r.status === grpc.StatusOK,
    });

    // Cannot get a pipeline of a non-exist user
    check(client.invoke('vdp.pipeline.v1alpha.PipelinePublicService/GetPipeline', {
      name: `pipelines/${reqBody.id}`
    }, constant.paramsWithJwt), {
      [`[with random "jwt-sub" header] vdp.pipeline.v1alpha.PipelinePublicService/GetPipeline response StatusUnknown`]: (r) => r.status === grpc.StatusUnknown,
    })

    // Delete the pipeline
    check(client.invoke(`vdp.pipeline.v1alpha.PipelinePublicService/DeletePipeline`, {
      name: `pipelines/${reqBody.id}`
    }), {
      [`vdp.pipeline.v1alpha.PipelinePublicService/DeletePipeline response StatusOK`]: (r) => r.status === grpc.StatusOK,
    });

    client.close();
  });
}

export function CheckUpdate() {

  group(`Pipelines API: Update a pipeline [with random "jwt-sub" header]`, () => {

    client.connect(constant.pipelineGRPCPublicHost, {
      plaintext: true
    });

    var reqBody = Object.assign({
      id: randomString(10),
    },
      constant.detSyncHTTPSingleModelRecipe
    )

    // Create a pipeline
    var resOrigin = client.invoke('vdp.pipeline.v1alpha.PipelinePublicService/CreatePipeline', {
      pipeline: reqBody
    })

    check(resOrigin, {
      [`vdp.pipeline.v1alpha.PipelinePublicService/CreatePipeline response StatusOK`]: (r) => r.status === grpc.StatusOK,
    });

    var reqBodyUpdate = Object.assign({
      id: reqBody.id,
      name: `pipelines/${reqBody.id}`,
      uid: "output-only-to-be-ignored",
      mode: "MODE_ASYNC",
      description: randomString(50),
    },)

    // Cannot update a pipeline of a non-exist user
    check(client.invoke('vdp.pipeline.v1alpha.PipelinePublicService/UpdatePipeline', {
      pipeline: reqBodyUpdate,
      update_mask: "description"
    }, constant.paramsWithJwt), {
      [`[with random "jwt-sub" header] vdp.pipeline.v1alpha.PipelinePublicService/UpdatePipeline response StatusUnknown`]: (r) => r.status === grpc.StatusUnknown,
    })

    // Delete the pipeline
    check(client.invoke(`vdp.pipeline.v1alpha.PipelinePublicService/DeletePipeline`, {
      name: `pipelines/${reqBody.id}`
    }), {
      [`vdp.pipeline.v1alpha.PipelinePublicService/DeletePipeline response StatusOK`]: (r) => r.status === grpc.StatusOK,
    });

    client.close()
  });
}

export function CheckUpdateState() {

  group(`Pipelines API: Update a pipeline state [with random "jwt-sub" header]`, () => {

    client.connect(constant.pipelineGRPCPublicHost, {
      plaintext: true
    });

    var reqBodySync = Object.assign({
      id: randomString(10),
    },
      constant.detSyncHTTPSingleModelRecipe
    )

    check(client.invoke('vdp.pipeline.v1alpha.PipelinePublicService/CreatePipeline', {
      pipeline: reqBodySync
    }), {
      [`vdp.pipeline.v1alpha.PipelinePublicService/CreatePipeline Sync response StatusOK`]: (r) => r.status === grpc.StatusOK,
      [`vdp.pipeline.v1alpha.PipelinePublicService/CreatePipeline Sync response pipeline state ACTIVE`]: (r) => r.message.pipeline.state === "STATE_ACTIVE",
    })

    // Cannot activate a pipeline of a non-exist user
    check(client.invoke('vdp.pipeline.v1alpha.PipelinePublicService/ActivatePipeline', {
      name: `pipelines/${reqBodySync.id}`
    }, constant.paramsWithJwt), {
      [`[with random "jwt-sub" header] vdp.pipeline.v1alpha.PipelinePublicService/ActivatePipeline response StatusUnknown`]: (r) => r.status === grpc.StatusUnknown,
    })

    // Cannot deactivate a pipeline of a non-exist user
    check(client.invoke('vdp.pipeline.v1alpha.PipelinePublicService/DeactivatePipeline', {
      name: `pipelines/${reqBodySync.id}`
    }, constant.paramsWithJwt), {
      [`[with random "jwt-sub" header] vdp.pipeline.v1alpha.PipelinePublicService/DeactivatePipeline response StatusUnknown`]: (r) => r.status === grpc.StatusUnknown,
    })

    // Delete the pipeline
    check(client.invoke(`vdp.pipeline.v1alpha.PipelinePublicService/DeletePipeline`, {
      name: `pipelines/${reqBodySync.id}`
    }), {
      [`vdp.pipeline.v1alpha.PipelinePublicService/DeletePipeline response StatusOK`]: (r) => r.status === grpc.StatusOK,
    });

    client.close()
  });
}

export function CheckRename() {

  group(`Pipelines API: Rename a pipeline [with random "jwt-sub" header]`, () => {

    client.connect(constant.pipelineGRPCPublicHost, {
      plaintext: true
    });

    var reqBody = Object.assign({
      id: randomString(10),
    },
      constant.detSyncHTTPSingleModelRecipe
    )

    // Create a pipeline
    var res = client.invoke('vdp.pipeline.v1alpha.PipelinePublicService/CreatePipeline', {
      pipeline: reqBody
    })

    check(res, {
      [`vdp.pipeline.v1alpha.PipelinePublicService/CreatePipeline response StatusOK`]: (r) => r.status === grpc.StatusOK,
      [`vdp.pipeline.v1alpha.PipelinePublicService/CreatePipeline response pipeline name`]: (r) => r.message.pipeline.name === `pipelines/${reqBody.id}`,
    });

    var new_pipeline_id = randomString(10)

    // Cannot rename a pipeline of a non-exist user
    check(client.invoke('vdp.pipeline.v1alpha.PipelinePublicService/RenamePipeline', {
      name: `pipelines/${reqBody.id}`,
      new_pipeline_id: new_pipeline_id
    }, constant.paramsWithJwt), {
      [`[with random "jwt-sub" header] vdp.pipeline.v1alpha.PipelinePublicService/RenamePipeline response StatusUnknown`]: (r) => r.status === grpc.StatusUnknown,
    })

    // Delete the pipeline
    check(client.invoke(`vdp.pipeline.v1alpha.PipelinePublicService/DeletePipeline`, {
      name: `pipelines/${reqBody.id}`
    }), {
      [`vdp.pipeline.v1alpha.PipelinePublicService/DeletePipeline response StatusOK`]: (r) => r.status === grpc.StatusOK,
    });

    client.close()
  });

}

export function CheckLookUp() {

  group(`Pipelines API: Look up a pipeline by uid [with random "jwt-sub" header]`, () => {

    client.connect(constant.pipelineGRPCPublicHost, {
      plaintext: true
    });

    var reqBody = Object.assign({
      id: randomString(10),
    },
      constant.detSyncHTTPSingleModelRecipe
    )

    // Create a pipeline
    var res = client.invoke('vdp.pipeline.v1alpha.PipelinePublicService/CreatePipeline', {
      pipeline: reqBody
    })

    check(res, {
      [`vdp.pipeline.v1alpha.PipelinePublicService/CreatePipeline response StatusOK`]: (r) => r.status === grpc.StatusOK,
    });

    // Cannot look up a pipeline of a non-exist user
    check(client.invoke('vdp.pipeline.v1alpha.PipelinePublicService/LookUpPipeline', {
      permalink: `pipelines/${res.message.pipeline.uid}`
    }, constant.paramsWithJwt), {
      [`[with random "jwt-sub" header] vdp.pipeline.v1alpha.PipelinePublicService/LookUpPipeline response StatusUnknown`]: (r) => r.status === grpc.StatusUnknown,
    })

    // Delete the pipeline
    check(client.invoke(`vdp.pipeline.v1alpha.PipelinePublicService/DeletePipeline`, {
      name: `pipelines/${reqBody.id}`
    }), {
      [`vdp.pipeline.v1alpha.PipelinePublicService/DeletePipeline response StatusOK`]: (r) => r.status === grpc.StatusOK,
    });

    client.close()
  });

}