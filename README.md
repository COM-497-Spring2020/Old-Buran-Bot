# Old-Buran-Bot
Group 1's Discord Bot, Shane Irons and Jordan Mason.

This bot should allow users of the official Monster Sanctuary Discord to input and later call their high scores in both the PvP leaderboard and the Infinity Arena. Additional features include the ability to remember one's PvP team.


```+-------------+------------+------+-----+---------------------+-------------------------------+
| Field       | Type       | Null | Key | Default             | Extra                         |
+-------------+------------+------+-----+---------------------+-------------------------------+
| DiscordID   | bigint(18) | YES  |     | NULL                |                               |
| RatingType  | bit(1)     | YES  |     | NULL                |                               |
| RatingScore | int(3)     | YES  |     | NULL                |                               |
| RatingImage | bit(1)     | YES  |     | NULL                |                               |
| TimeStamp   | timestamp  | NO   |     | 0000-00-00 00:00:00 | on update current_timestamp() |
+-------------+------------+------+-----+---------------------+-------------------------------+```
