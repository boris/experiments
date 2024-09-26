use rand::seq::SliceRandom;

fn main() {
    let vs = vec!["foo", "bar", "baz"];
    println!("{:?}", vs.choose(&mut rand::thread_rng()));
}

