from telegram.ext import Updater, InlineQueryHandler, CommandHandler
import requests
import re

TOKEN = ''
def start(update, context):
    print(update.message.from_user)
    user = update.message.from_user['username']
    """Send a message when the command /start is issued."""
    update.message.reply_text('Hi! @{}'.format(user))


def main():
    """Start the bot."""
    # Create the Updater and pass it your bot's token.
    # Make sure to set use_context=True to use the new context based callbacks
    # Post version 12 this will no longer be necessary
    updater = Updater(TOKEN, use_context=True)

    # Get the dispatcher to register handlers
    dp = updater.dispatcher

    # on different commands - answer in Telegram
    dp.add_handler(CommandHandler("start", start))

    # Start the Bot
    updater.start_polling()

    # Run the bot until you press Ctrl-C or the process receives SIGINT,
    # SIGTERM or SIGABRT. This should be used most of the time, since
    # start_polling() is non-blocking and will stop the bot gracefully.
    updater.idle()

def price():
    r = requests.get('https://www.huobi.com/en-us/exchange/eth_usdt/').data()
    print(r)

if __name__ == '__main__':
    price()