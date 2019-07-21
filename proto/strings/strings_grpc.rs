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

const METHOD_STRINGS_EXTRACTOR_EXTRACT_STRINGS: ::grpcio::Method<super::strings::StringsRequest, super::strings::StringsReply> = ::grpcio::Method {
    ty: ::grpcio::MethodType::Unary,
    name: "/strings.StringsExtractor/ExtractStrings",
    req_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
    resp_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
};

#[derive(Clone)]
pub struct StringsExtractorClient {
    client: ::grpcio::Client,
}

impl StringsExtractorClient {
    pub fn new(channel: ::grpcio::Channel) -> Self {
        StringsExtractorClient {
            client: ::grpcio::Client::new(channel),
        }
    }

    pub fn extract_strings_opt(&self, req: &super::strings::StringsRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<super::strings::StringsReply> {
        self.client.unary_call(&METHOD_STRINGS_EXTRACTOR_EXTRACT_STRINGS, req, opt)
    }

    pub fn extract_strings(&self, req: &super::strings::StringsRequest) -> ::grpcio::Result<super::strings::StringsReply> {
        self.extract_strings_opt(req, ::grpcio::CallOption::default())
    }

    pub fn extract_strings_async_opt(&self, req: &super::strings::StringsRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::strings::StringsReply>> {
        self.client.unary_call_async(&METHOD_STRINGS_EXTRACTOR_EXTRACT_STRINGS, req, opt)
    }

    pub fn extract_strings_async(&self, req: &super::strings::StringsRequest) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::strings::StringsReply>> {
        self.extract_strings_async_opt(req, ::grpcio::CallOption::default())
    }
    pub fn spawn<F>(&self, f: F) where F: ::futures::Future<Item = (), Error = ()> + Send + 'static {
        self.client.spawn(f)
    }
}

pub trait StringsExtractor {
    fn extract_strings(&mut self, ctx: ::grpcio::RpcContext, req: super::strings::StringsRequest, sink: ::grpcio::UnarySink<super::strings::StringsReply>);
}

pub fn create_strings_extractor<S: StringsExtractor + Send + Clone + 'static>(s: S) -> ::grpcio::Service {
    let mut builder = ::grpcio::ServiceBuilder::new();
    let mut instance = s.clone();
    builder = builder.add_unary_handler(&METHOD_STRINGS_EXTRACTOR_EXTRACT_STRINGS, move |ctx, req, resp| {
        instance.extract_strings(ctx, req, resp)
    });
    builder.build()
}
