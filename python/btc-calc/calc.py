import click

@click.command()
@click.option('--flo', default=15.0, type=float, help="Flo's share. Default: 15.0")
@click.option('--paz', default=40.0, type=float, help="Paz's share. Default: 40.0")
@click.option('--boris', default=40.0, type=float, help="Boris' share. Default: 40.0")
@click.argument('btc_price', type=float)
def calc_price(btc_price, flo, paz, boris):
    """
    Calculate the amount of BTC on each weekly buy. It uses BTC_PRICE to set the value of BTC in USD.

    Example:

        python calc.py --flo 10 --paz 150 --boris 21 50000

    """
    click.echo(f'Flo:\t {flo / btc_price:.8f} BTC')
    click.echo(f'Paz:\t {paz / btc_price:.8f} BTC')
    click.echo(f'Boris:\t {boris / btc_price:.8f} BTC')

if __name__ == '__main__':
    calc_price()
