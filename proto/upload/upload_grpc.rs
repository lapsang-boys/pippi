// This file is generated. Do not edit
// @generated

// https://github.com/Manishearth/rust-clippy/issues/702
#![allow(unknown_lints)]
#![allow(clippy::all)]

#![cfg_attr(rustfmt, rustfmt_skip)]

#![allow(box_pointers)]
#![allow(dead_code)]
#![allow(missing_docs)]
#![allow(non_camel_case_types)]
#![allow(non_snake_case)]
#![allow(non_upper_case_globals)]
#![allow(trivial_casts)]
#![allow(unsafe_code)]
#![allow(unused_imports)]
#![allow(unused_results)]

const METHOD_UPLOAD_UPLOAD: ::grpcio::Method<super::upload::UploadRequest, super::upload::UploadReply> = ::grpcio::Method {
    ty: ::grpcio::MethodType::Unary,
    name: "/upload.Upload/Upload",
    req_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
    resp_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
};

#[derive(Clone)]
pub struct UploadClient {
    client: ::grpcio::Client,
}

impl UploadClient {
    pub fn new(channel: ::grpcio::Channel) -> Self {
        UploadClient {
            client: ::grpcio::Client::new(channel),
        }
    }

    pub fn upload_opt(&self, req: &super::upload::UploadRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<super::upload::UploadReply> {
        self.client.unary_call(&METHOD_UPLOAD_UPLOAD, req, opt)
    }

    pub fn upload(&self, req: &super::upload::UploadRequest) -> ::grpcio::Result<super::upload::UploadReply> {
        self.upload_opt(req, ::grpcio::CallOption::default())
    }

    pub fn upload_async_opt(&self, req: &super::upload::UploadRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::upload::UploadReply>> {
        self.client.unary_call_async(&METHOD_UPLOAD_UPLOAD, req, opt)
    }

    pub fn upload_async(&self, req: &super::upload::UploadRequest) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::upload::UploadReply>> {
        self.upload_async_opt(req, ::grpcio::CallOption::default())
    }
    pub fn spawn<F>(&self, f: F) where F: ::futures::Future<Item = (), Error = ()> + Send + 'static {
        self.client.spawn(f)
    }
}

pub trait Upload {
    fn upload(&mut self, ctx: ::grpcio::RpcContext, req: super::upload::UploadRequest, sink: ::grpcio::UnarySink<super::upload::UploadReply>);
}

pub fn create_upload<S: Upload + Send + Clone + 'static>(s: S) -> ::grpcio::Service {
    let mut builder = ::grpcio::ServiceBuilder::new();
    let mut instance = s.clone();
    builder = builder.add_unary_handler(&METHOD_UPLOAD_UPLOAD, move |ctx, req, resp| {
        instance.upload(ctx, req, resp)
    });
    builder.build()
}
