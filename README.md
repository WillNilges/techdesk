# Cursed Tech Support

In Channel A on Slack A, user uses `/tech`, which opens up a dialogue. The dialogue asks for a subject and description.

On submission, the bot will post a message to Channel A on Slack A, and Channel B on Slack B. The message will be pinned on Channel B.

A Tech will service the ticket, and then unpin it. When it is unpinned, the Bot will note down the time it was unpinned, and update the message in Channel A with the time to resolution and ping user in the replies.

The biggest problem I see with this is storage. Databases. Whatever the fuck. How do I keep track of the tickets? I guess I can just use Slack as a database again, much to my fucking chagrin though.

So, I think whenever anything happens, we're going to need to pull up a message. I think I can use the timestamp of the message that has been interacted with as a sort of "anchor" and then just search a range (there's going to be some delay since the event has to go to the server and back to Slack).

So, I probably won't need the channel history. I can just look up in Slack using the timestamp from the message on the event.

Maybe I should make a TechChannel and AdminChannel object. Or maybe I should have one object that has both sockets

How about we canonize the "Tech" and "Admin" roles? "Tech" is Marco and "Admin" is Allen or whoever that guy is.

