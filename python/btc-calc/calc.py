import click

@click.command()
@click.option('--flo', default=15.0, type=float)
@click.option('--paz', default=40.0, type=float)
@click.option('--boris', default=40.0, type=float)
@click.argument('btc_price', type=float)
def calc_price(btc_price, flo, paz, boris):
    click.echo(f'Flo:\t {flo / btc_price:.8f} BTC')
    click.echo(f'Paz:\t {paz / btc_price:.8f} BTC')
    click.echo(f'Boris:\t {boris / btc_price:.8f} BTC')

if __name__ == '__main__':
    calc_price()
