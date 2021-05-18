# bot.py

import os
from datetime import datetime
from bs4 import BeautifulSoup
import requests
import discord
from discord.ext import commands
# from dotenv import load_dotenv

# load_dotenv()
# BOT_TOKEN = os.getenv('DISCORD_TOKEN')
#todo use env for deployed versions
BOT_TOKEN = 'NzkwMDI3MjAwMzUxMDQzNjI0.X96oKg.qExzgc5gZQvG05EljgwyDeFinI8'

is_marathon = False
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
async def manual_update(context):
    global is_marathon
    global last_offer

    offer_info = get_latest_lbw_offer()
    new_offer = is_new_offer(offer_info)
    marathon_change = is_marathon == offer_info['is_marathon']

    last_offer = offer_info
    is_marathon = offer_info['is_marathon']
    
    if is_marathon: # todo notify if marathon has changed in state
        await message.channel.send('marathon has started')

    await context.send(embed=create_notification_embed(new_offer, offer_info))

@bot.command(name='is-marathon', help='checks whether a marathon is currently underway.')
async def start_interval(context):
    await context.send('yes' if is_marathon else 'no')

@bot.command(name='start', help='starts the notifier scheduling')
async def start_interval(context):
    await context.send('test')

@bot.command(name='stop', help='stops the notifier scheduling')
async def stop_interval(context):
    await context.send('test')

@bot.command(name='set-interval', help='sets the notifier scheduling interval in minutes. if `auto`, itr will set it to automatically respond per day')
async def set_interval(context, arg):
    await context.send(f'test {arg}')

@bot.command(name='mute', help='disables the notifier offer posts of the specified type (`all`, `offers`, or `marathon`)')
@commands.has_role('admin')
async def mute_notifications(context, notification_type):
    if notification_type == 'all':
        # mute all notifications
        pass
    elif notification_type == 'offers': 
        # mute jsut offers
        pass
    elif notification_type == 'marathon': 
        # mute jsut marathon
        pass
    else:
        await on_command_error(context, 'Unrecognized notification type.\nPlease use `all`, `offers`, or `marathon`')
        return
    await context.send(f'test {notification_type}')
    

################################ Helper Functions


def get_latest_lbw_offer():
    soup = BeautifulSoup(requests.get('https://www.lastbottlewines.com/').text, 'lxml')
    
    is_marathon = soup.find('div', attrs={ 'class': 'marquee-top' }) or len(soup.find_all('div', attrs={ 'class': 'marathon' })) > 0

    offer_name = soup.find('h1', attrs={ 'class': 'offer-name' }).text
    offer_image = 'http:' + soup.find('img', id='offer-image').attrs['src']

    price_containers = soup.find_all('div', attrs={ 'class': 'price-holder' }, limit=3)
    price_data = list(map(lambda container: (container.find('p').text, container.find('span', attrs={ 'class': 'amount' }).text), price_containers))

    return {
        'is_marathon': is_marathon,
        'name': offer_name,
        'image': offer_image,
        'prices': price_data
    }

def create_notification_embed(new_offer, offer_info):
    global is_marathon
    # todo, if new offer , color white

    embed = discord.Embed(title=f':wine_glass:{"New" if new_offer else "Current"} Offer:wine_glass:', colour=discord.Colour(0xf03d44) if is_marathon else 0, url="https://www.lastbottlewines.com/", timestamp=datetime.now())

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