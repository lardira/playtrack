### GET /api/leaderboard

        {
        "users": [
            {
            "place": 1,
            "user_id": 1,
            "nickname": "Walun",
            "current_game": "Dark Souls 3 Pojarnie",
            "points": 2,
            "completed": 4,
            "dropped": 12,
            "reroll": 1
            }
        ],
        "stats": {
            "week": {
            "points": [120, 140, 180, 200, 220, 300, 380],
            "completed": [1, 0, 2, 1, 3, 2, 3]
            },
            "month": { ... }
        }
        }

---

### GET /api/users/:id

    {
    "user": {
        "id": 1,
        "nickname": "Walun",
        "avatar": "/avatars/smeiP3tbc4rKm.png",
        "status": "drop",
        "about": "Это моя",
        "stats": {
        "points": 4,
        "completed": 4,
        "dropped": 12,
        "reroll": 1
        },
        "is_owner": true
    },
    "games": [
        {
        "date": "2025-01-10",
        "time_spent": "02:35",
        "game": {
            "title": "Dark Souls 3",
            "url": "https://store.steampowered.com/..."
        },
        "status": "Пройдено",
        "comment": "Боль, но кайф",
        "rating": "99/100"
        }
    ]
    }

---

POST /api/games - 
POST /api/login - JWT или ставит cookie
