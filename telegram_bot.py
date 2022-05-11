import os
from telegram.ext import Updater, InlineQueryHandler, CommandHandler
from pycoingecko import CoinGeckoAPI
from dotenv import load_dotenv
import json
from price_list import COIN_MAP

load_dotenv()
TOKEN = os.getenv("TOKEN")

MESSAGE = '''
Hi, @{}!
Here are the commands instructions:
(Only support coin/usd format)
- /set lo 'coin' price
    Alert when price fell below target price
    eg: /set lo eth 2100

- /set hi 'coin' price
    Alert when price raised above target price
    eg: /set hi eth 2100
  
- /del lo 'coin'
    Delete lower coin alert
    eg: /del lo eth 
  
- /del hi 'coin'
    Delete higher coin alert
    eg: /del hi eth
    
- /alert
    Check on alert list
    
- /price 'coin'
    Check on coin price
    eg: /price eth
'''

def help(update, context):
    user = update.message.from_user['username']
    update.message.reply_text(MESSAGE.format(user))

def check_price(update, context):
    cg = CoinGeckoAPI()
    coin = update.message.text.split(" ")[1]
    id = read_from_price_list(coin)

    if id is None:
        update.message.reply_text("No such coin named: {}".format(coin))
    else:
        prices = get_price(cg, id)
        update.message.reply_text("{} current price: U$ {}".format(coin.upper(), prices))

def main():
    """Start the bot."""
    # Create the Updater and pass it your bot's token.
    updater = Updater(TOKEN, use_context=True)

    # Get the dispatcher to register handlers
    dp = updater.dispatcher

    # Command List
    dp.add_handler(CommandHandler("help", help))
    dp.add_handler(CommandHandler("price", check_price))

    # Start the Bot
    updater.start_polling()

    # Run the bot until you press Ctrl-C or the process receives SIGINT,
    updater.idle()

def update_price_list(cg):
    price_list = cg.get_coins_list()
    coin_map = {}
    for item in price_list:
        coin_map[item['symbol'].lower()] = item['id']

    s = json.dumps(coin_map)
    f = open("price_list.py", "w")
    f.write("COIN_MAP = {}".format(s))
    f.close()

def read_from_price_list(coin_sym):
    coin_sym = coin_sym.lower()
    if coin_sym not in COIN_MAP:
        return None
    return COIN_MAP[coin_sym]

def get_price(cg, id):
    js = cg.get_price(ids=id, vs_currencies='usd')
    keys = list(js.keys())[0]
    price = js[keys]['usd']
    return ("%.10f" % price).rstrip('0').rstrip('.')
    
if __name__ == '__main__':
    # cg = CoinGeckoAPI()
    # id = read_from_price_list('caw')
    # print(get_price(cg, id))
    main()