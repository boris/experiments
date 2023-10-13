use serde::{Deserialize, Serialize};
use serde_yaml::{self};

#[derive(Debug, Serialize, Deserialize)]
struct Config {
    num_threads: u8,
    data_sources: Vec<String>,
    db: String,
}

fn main() {
    let f = std::fs::File::open("config.yml").expect("Unable to open file");
    let scrape_config: Config = serde_yaml::from_reader(f).expect("Unable to read yaml");

    println!("{:?}", scrape_config);
}
