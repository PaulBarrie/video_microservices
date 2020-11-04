SELECT * FROM video WHERE name LIKE :video_search;
SELECT * FROM user WHERE username LIKE :user_search OR pseudo LIKE :user_search;
