# 1. Connect to the shell

docker exec -it gi25-group2-projectmonorepo-mongodb-1 mongosh -u rootuser -p rootpass --authenticationDatabase admin

# 2. Once inside mongosh, run these commands:

use blog_db

# Paste the large insertMany command from above...

# 3. (Optional) test to see the insert

db.articles.find().pretty()

# 4. (Optional) drop collection

db.articles.deleteMany({})
exit

# Script

db.articles.insertMany([
{
title: "THE FUTURE OF QUANTUM COMPUTING",
content: "Quantum computing promises to revolutionize industries by solving complex problems beyond the reach of classical computers. This article explores the basics and potential impacts.",
publisherName: "Tech Today",
publisherId: 101,
category: "Technology",
createdAt: ISODate("2025-09-19T12:00:00Z"),
comments: [
{
publisherId: 23,
publisherName: "Alice",
content: "Incredibly insightful! The explanation of quantum superposition was the clearest I've ever read. Thank you!",
createdAt: ISODate("2025-09-20T10:00:00Z")
},
{
publisherId: 45,
publisherName: "Bob",
content: "Great overview. I'm excited to see how this will affect cryptography in the coming years.",
createdAt: ISODate("2025-09-21T14:30:00Z")
}
]
},
{
title: "A GUIDE TO TRAVELING THROUGH SOUTHEAST ASIA",
content: "From the bustling streets of Bangkok to the serene beaches of Bali, Southeast Asia offers an adventure for every type of traveler. Here are the must-see destinations.",
publisherName: "Wanderlust Weekly",
publisherId: 205,
category: "Travel",
createdAt: ISODate("2025-09-15T09:30:00Z"),
comments: [
{
publisherId: 17,
publisherName: "Charlie",
content: "This guide is fantastic! I just got back from Vietnam and this article perfectly captures the experience.",
createdAt: ISODate("2025-09-16T11:20:00Z")
},
{
publisherId: 88,
publisherName: "Diana",
content: "Are there any specific budget tips for traveling through the Philippines? Planning a trip for next month!",
createdAt: ISODate("2025-09-17T18:05:00Z")
}
]
},
{
title: "SIMPLE TIPS FOR A HEALTHIER LIFESTYLE",
content: "Improving your health doesn't have to be complicated. Discover five easy-to-implement habits that can make a significant difference in your daily well-being.",
publisherName: "Wellness Hub",
publisherId: 310,
category: "Health",
createdAt: ISODate("2025-09-10T15:00:00Z"),
comments: [
{
publisherId: 5,
publisherName: "Eve",
content: "Wonderful post! The tip about mindful eating has already made a huge difference for me.",
createdAt: ISODate("2025-09-11T08:00:00Z")
},
{
publisherId: 62,
publisherName: "Frank",
content: "Solid advice. Consistency really is the key. Thanks for the reminder.",
createdAt: ISODate("2025-09-12T16:45:00Z")
}
]
},
{
title: "MASTERING THE ART OF SOURDOUGH",
content: "Baking the perfect loaf of sourdough bread is a rewarding experience. This guide covers everything from creating your starter to achieving the ideal crust.",
publisherName: "The Home Baker",
publisherId: 415,
category: "Food",
createdAt: ISODate("2025-08-28T18:45:00Z"),
comments: [
{
publisherId: 71,
publisherName: "Grace",
content: "My starter never seems to be active enough. Any tips for a colder kitchen?",
createdAt: ISODate("2025-08-29T09:00:00Z")
},
{
publisherId: 99,
publisherName: "Heidi",
content: "This is the guide I wish I had when I started baking. The section on scoring is particularly helpful!",
createdAt: ISODate("2025-08-30T12:15:00Z")
}
]
},
{
title: "AI IN 2025: WHAT'S NEW?",
content: "Artificial intelligence continues to evolve at a rapid pace. We'll look at the latest breakthroughs and what they mean for the future of technology and society.",
publisherName: "Tech Today",
publisherId: 101,
category: "Technology",
createdAt: ISODate("2025-08-25T11:20:00Z"),
comments: [
{
publisherId: 10,
publisherName: "Ivan",
content: "The progress in generative models is both exciting and a bit scary. The ethical implications need more discussion.",
createdAt: ISODate("2025-08-26T13:00:00Z")
},
{
publisherId: 22,
publisherName: "Judy",
content: "I'm most interested in how this will affect personalized medicine. The potential is enormous.",
createdAt: ISODate("2025-08-27T10:55:00Z")
}
]
}
]);
