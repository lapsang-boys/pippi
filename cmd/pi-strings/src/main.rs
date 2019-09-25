#[macro_use]
extern crate log;

#[macro_use]
extern crate lazy_static;

use std::fs;
use std::str;

use std::io::Read;
use std::sync::Arc;
use std::{io, thread};

use futures::sync::oneshot;
use futures::Future;
use grpcio::{Environment, RpcContext, RpcStatus, RpcStatusCode, ServerBuilder, UnarySink};
use protobuf::RepeatedField;

#[path = "../../../proto/strings/strings.rs"]
mod strings;
#[path = "../../../proto/strings/strings_grpc.rs"]
mod strings_grpc;

use directories::BaseDirs;

use regex::Regex;

#[derive(Clone)]
struct StringsService;

fn is_valid(id: &str) -> bool {
    lazy_static! {
        static ref RE: Regex = Regex::new("^[0-9a-f]{64}$").unwrap();
    }

    RE.is_match(id) && id.len() == 64
}

impl strings_grpc::StringsExtractor for StringsService {
    fn extract_strings(
        &mut self,
        ctx: RpcContext<'_>,
        req: strings::StringsRequest,
        sink: UnarySink<strings::StringsReply>,
    ) {
        let id = req.get_id();
        println!("{}", id);
        println!("{}", is_valid(id));
        if !is_valid(id) {
            let status = RpcStatus::new(
                RpcStatusCode::InvalidArgument,
                Some("Invalid argument".to_string()),
            );
            let f = sink
                .fail(status)
                .map_err(move |e| error!("failed to reply {:?}: {:?}", req, e));
            ctx.spawn(f);
            return;
        }
        let base_dirs = BaseDirs::new().unwrap();
        let filename = base_dirs
            .cache_dir()
            .join("pippi")
            .join(id)
            .join(id.to_owned() + ".bin")
            .to_str()
            .unwrap()
            .to_string();

        println!("{}", filename);
        let sinfos = extract_strings_from_path(filename).unwrap();
        let mut resp = strings::StringsReply::default();
        resp.set_strings(RepeatedField::from_vec(sinfos));
        let f = sink
            .success(resp)
            .map_err(move |e| error!("failed to reply {:?}: {:?}", req, e));
        ctx.spawn(f)
    }
}

fn main() {
    let env = Arc::new(Environment::new(1));
    let service = strings_grpc::create_strings_extractor(StringsService);
    let mut server = ServerBuilder::new(env)
        .register_service(service)
        .bind("127.0.0.1", 1400)
        .build()
        .unwrap();
    server.start();
    for &(ref host, port) in server.bind_addrs() {
        println!("listening on {}:{}", host, port);
    }
    let (tx, rx) = oneshot::channel();
    //thread::spawn(move || {
    //    info!("Press ENTER to exit...");
    //    let _ = io::stdin().read(&mut [0]).unwrap();
    //    tx.send(())
    //});
    if (false) {
        tx.send(());
    }
    let _ = rx.wait();
    let _ = server.shutdown().wait();
}

const MIN_LENGTH: usize = 4;

fn is_ascii_printable(c: u8) -> bool {
    return c.is_ascii_graphic() || c == ' ' as u8 || c == '\t' as u8;
}

fn extract_strings_from_path(filename: String) -> Result<Vec<strings::StringInfo>, std::io::Error> {
    let contents = fs::read(filename)?;

    let mut a: Vec<u8> = Vec::new();
    let mut sinfos: Vec<strings::StringInfo> = Vec::new();
    for (index, c) in contents.iter().enumerate() {
        if is_ascii_printable(*c) {
            a.push(*c);
            continue;
        } else {
            if a.len() < MIN_LENGTH {
                a.clear();
                continue;
            }
            let inp = a.clone();
            let s = match str::from_utf8(&inp) {
                Ok(s) => s,
                Err(_) => "",
            };
            if s.len() > 0 {
                let mut sinfo = strings::StringInfo::default();
                sinfo.set_location((index - s.len()) as u64);
                sinfo.set_raw_string(s.to_string());
                sinfos.push(sinfo);
                debug!("{} {}", index - s.len(), s);
            }
            a.clear();
        }
    }

    Ok(sinfos)
}
