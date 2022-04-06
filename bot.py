# main.py

import os
from datetime import datetime
from bs4 import BeautifulSoup
import requests
import enum
import discord
from discord.ext import commands, tasks

class Notification(enum.Enum):
    All = 'all'
    Offers = 'offers'
    Marathon = 'marathon'

# for reference
NOTIFICATIONS_CHANNEL_ID = 829129887905742858
TESTING_CHANNEL_ID = 829063368756166667

# for local, every time do `export DISCORD_BOT_TOKEN=<value from token.txt>`
BOT_TOKEN = os.getenv('DISCORD_BOT_TOKEN') or ''
BOT_CHANNEL_ID = os.getenv('CHANNEL_ID') or TESTING_CHANNEL_ID

is_marathon = False
notification_rule = Notification.All
last_offer = {}

bot = commands.Bot(command_prefix='!')

################################ Lifecycle Hooks

@bot.event
async def on_ready():
    #todo start up interval now
    print(f'{bot.user.name} has connected to Discord!')

@bot.event
async def on_command_error(context, error):
    await context.send(f'An error occurred: {error}')


################################ Commands


@bot.command(name='update', help='gets the current LBW offer')
async def update_offer(context):
    is_new_offer, is_new_marathon = update_stored_offer_info()

    if is_new_marathon:
        await context.send('!!!!!!!!!!!!!!!!! Marathon has started !!!!!!!!!!!!!!!!!')

    await context.send(embed=create_notification_embed(last_offer, is_new_offer))

@bot.command(name='is-marathon', help='checks whether a marathon is currently underway.')
async def start_interval(context):
    await context.send('yes' if is_marathon else 'no')

@bot.command(name='start', help='starts the notifier scheduling')
async def start_notification_interval(context, interval):
    try:
        interval_number = int(interval)
    except:
        await context.send('Unrecognized interval. Please enter a number for the interval')
        return

    change_task_interval(minutes=interval_number)
    check_offer_and_notify.start()

    await context.send(f"Notifications on an interval of {interval_number} minutes started.\nType command 'stop' to stop the interval")

@bot.command(name='stop', help='stops the notifier scheduling')
async def stop_notification_interval(context, hard_stop = False):
    if hard_stop:
        check_offer_and_notify.cancel()
    else:
        check_offer_and_notify.stop()
    
    await context.send('Notificaiton interval stopped')

@bot.command(name='set-interval', help='sets the notifier scheduling interval in minutes or hours. Use format <n> <m=minutes | h=hours> ')
async def set_interval(context, time, span):
    span_string = ''
    if span == 'h':
        span_string = 'hours' if time != 1 else 'hour'
        change_task_interval(hours=int(time))
    elif span == 'm':
        span_string = 'minutes' if time != 1 else 'minute'
        change_task_interval(minutes=int(time))
    elif span == 's':
        span_string = 'seconds' if time != 1 else 'second'
        change_task_interval(seconds=int(time))
    else:
        await context.send(f"Unrecognized interval notation. Please use 'm' for minutes or 'h' for hours")
        return
    
    await context.send(f'setting interval to {time} {span_string}')

@bot.command(name='mute', help='changes the posts of the bot to only notify for the specified type (`all`, `offers`, or `marathon`)')
@commands.has_role('admin')
async def change_notification_rule(context, notification_type):
    global notification_rule

    if notification_type != Notification.All or notification_type != Notification.Offers or notification_type != Notification.Marathon:
        await on_command_error(context, "Unrecognized notification type.\nPlease use 'all', 'offers', or 'marathon'")
        return

    notification_rule = notification_type

    await context.send(f'changing notifications to only show for {notification_type}')
    

################################ Tasks


@tasks.loop(minutes=1)
async def check_offer_and_notify():
    channel = bot.get_channel(BOT_CHANNEL_ID)
    
    is_new_offer, is_new_marathon = update_stored_offer_info()
    
    if is_new_marathon and (notification_rule == Notification.All or notification_rule == Notification.Marathon):
        await channel.send('!!!!!!!!!!!!!!!!! Marathon has started !!!!!!!!!!!!!!!!!')

    if is_new_offer:
        can_send_marathon = is_marathon and notification_rule == Notification.Marathon
        can_send_offer = (not is_marathon) and notification_rule == Notification.Offers
        if notification_rule == Notification.All or can_send_marathon or can_send_offer:
            await channel.send(embed=create_notification_embed(last_offer, True))


################################ Helper Functions

def change_task_auto_interval():
    if is_marathon:
        check_offer_and_notify.change_interval(minutes=1)
    else:
        check_offer_and_notify.change_interval(hours=24)

def change_task_interval(*, seconds=0, minutes = 0, hours = 0):
    check_offer_and_notify.change_interval(seconds=seconds, minutes=minutes, hours=hours)

def get_latest_offer():
    soup = BeautifulSoup(requests.get('https://www.lastbottlewines.com/').text, 'lxml')
    
    is_marathon = soup.find('div', attrs={ 'class': 'marquee-top' }) or len(soup.find_all('div', attrs={ 'class': 'marathon' })) > 0

    offer_name = soup.find('h1', attrs={ 'class': 'offer-name' }).text
    image_src = soup.find('img', id='offer-image')
    offer_image = ''
    if not image_src is None:
        offer_image = 'http:' + soup.find('img', id='offer-image').attrs['src']

    price_containers = soup.find_all('div', attrs={ 'class': 'price-holder' }, limit=3)
    price_data = list(map(lambda container: (container.find('p').text, container.find('span', attrs={ 'class': 'amount' }).text), price_containers))

    return {
        'is_marathon': is_marathon,
        'name': offer_name,
        'image': offer_image,
        'prices': price_data
    }

def update_stored_offer_info():
    global is_marathon
    global last_offer

    offer_info = get_latest_offer()
    new_offer = is_new_offer(offer_info)
    marathon_change = is_marathon != offer_info['is_marathon']

    last_offer = offer_info
    is_marathon = offer_info['is_marathon']
    
    return new_offer, marathon_change and is_marathon

def create_notification_embed(offer_info, is_new_offer):
    global is_marathon
    # todo, if new offer , color white

    embed = discord.Embed(title=f':wine_glass:{"New" if is_new_offer else "Current"} Offer:wine_glass:', colour=discord.Colour(0xf03d44) if is_marathon else 0, url="https://www.lastbottlewines.com/")

    embed.set_image(url=offer_info['image'])
    embed.set_thumbnail(url="https://www.lastbottlewines.com/favicon.png")

    embed.add_field(name="Name", value=offer_info['name'], inline=False)

    # todo put discout in there as well
    for price_data in offer_info["prices"]:
        embed.add_field(name=price_data[0], value=f'${price_data[1]}', inline=True)
    
    tokenized_name = offer_info['name'].replace(' ', '+')
    url_tokenized_name = offer_info['name'].replace(' ', '%20')
    
    google_search_link = f'https://www.google.com/search?q={tokenized_name}&oq={tokenized_name}&aqs=chrome..69i57j69i61.1483j0j4&sourceid=chrome&ie=UTF-8'
    vivino_search_link = f'https://www.vivino.com/search/wines?q={tokenized_name}'
    wine_searcher_search_link = f'https://www.wine-searcher.com/find/{tokenized_name}'
    cellar_tracker_search_link = f'https://www.cellartracker.com/list.asp?fInStock=0&Table=List&iUserOverride=0&szSearch={tokenized_name}#selected%3DW3898544_1_Kcd91961cd0541a38650e2d05f7aa1b3f'

    binnys_search_link = f'https://www.binnys.com/search?q={url_tokenized_name}'
    vin_search_link = f'https://vinchicago.com/wines/search?keyword={tokenized_name}&limitstart=0&option=com_virtuemart&view=category'

    embed.add_field(name='Search Links', value=f'[Google]({google_search_link})\n[Vivino]({vivino_search_link})\n[Wine Searcher]({wine_searcher_search_link})\n[Cellar Tracker]({cellar_tracker_search_link})', inline=False)
    embed.add_field(name='Shop Search Links', value=f'[Binny\'s]({binnys_search_link})\n[Vin]({vin_search_link})', inline=False)

    return embed

def is_new_offer(offer):
    global last_offer

    if not 'name' in last_offer: # no last offer
        return True

    return not last_offer['name'] == offer['name']


bot.run(BOT_TOKEN)