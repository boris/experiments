use reqwest;
use serde::Deserialize;
use serde_json::Value;

#[derive(Deserialize, Debug)]
struct BpiResponse {
    bpi: BpiData,
}

#[derive(Deserialize, Debug)]
struct BpiData {
    #[serde(rename = "USD")]
    usd: Option<CurrencyData>,
}

#[derive(Deserialize, Debug)]
struct CurrencyData {
    code: String,
    rate_float: f64,
}

fn main() {
    let url = "https://api.coindesk.com/v1/bpi/currentprice/USD.json";
    let response = reqwest::blocking::get(url).expect("Failed to send request");

    if response.status().is_success() {
        let json: Value = serde_json::from_str(&response.text().expect("Failed to read response")).expect("Failed to parse JSON");
        let bpi_response: BpiResponse = serde_json::from_value(json).expect("Failed to deserialize JSON");

        if let Some(usd_data) = bpi_response.bpi.usd {
            let usd_rate = usd_data.rate_float;
            let usd_description = usd_data.code;
            println!("{}: {}", usd_description, usd_rate);
        }


    } else {
        println!("Request failed with status code: {}", response.status());
    }
}
