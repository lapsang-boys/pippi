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

const METHOD_BINARY_PARSER_PARSE_BINARY: ::grpcio::Method<super::bin::ParseBinaryRequest, super::bin::File> = ::grpcio::Method {
    ty: ::grpcio::MethodType::Unary,
    name: "/bin.BinaryParser/ParseBinary",
    req_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
    resp_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
};

#[derive(Clone)]
pub struct BinaryParserClient {
    client: ::grpcio::Client,
}

impl BinaryParserClient {
    pub fn new(channel: ::grpcio::Channel) -> Self {
        BinaryParserClient {
            client: ::grpcio::Client::new(channel),
        }
    }

    pub fn parse_binary_opt(&self, req: &super::bin::ParseBinaryRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<super::bin::File> {
        self.client.unary_call(&METHOD_BINARY_PARSER_PARSE_BINARY, req, opt)
    }

    pub fn parse_binary(&self, req: &super::bin::ParseBinaryRequest) -> ::grpcio::Result<super::bin::File> {
        self.parse_binary_opt(req, ::grpcio::CallOption::default())
    }

    pub fn parse_binary_async_opt(&self, req: &super::bin::ParseBinaryRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::bin::File>> {
        self.client.unary_call_async(&METHOD_BINARY_PARSER_PARSE_BINARY, req, opt)
    }

    pub fn parse_binary_async(&self, req: &super::bin::ParseBinaryRequest) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::bin::File>> {
        self.parse_binary_async_opt(req, ::grpcio::CallOption::default())
    }
    pub fn spawn<F>(&self, f: F) where F: ::futures::Future<Item = (), Error = ()> + Send + 'static {
        self.client.spawn(f)
    }
}

pub trait BinaryParser {
    fn parse_binary(&mut self, ctx: ::grpcio::RpcContext, req: super::bin::ParseBinaryRequest, sink: ::grpcio::UnarySink<super::bin::File>);
}

pub fn create_binary_parser<S: BinaryParser + Send + Clone + 'static>(s: S) -> ::grpcio::Service {
    let mut builder = ::grpcio::ServiceBuilder::new();
    let mut instance = s.clone();
    builder = builder.add_unary_handler(&METHOD_BINARY_PARSER_PARSE_BINARY, move |ctx, req, resp| {
        instance.parse_binary(ctx, req, resp)
    });
    builder.build()
}
