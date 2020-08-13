# TestTaskGolang

Маршруты:
1) /getTokes 
  Читае GUID как параметр запроса, генерирует access и refresh токены для GUID и IP с которого отправлен запрос, сохраняет их в базу данных (mongoDB), возвращает токены в формате json с полями access_token и refresh_token
  
2)/refreshTokens 
  Принимает json с полями: _id (GUID), access_token, refresh_token, делает проверку refresh токена пользователя сравнивая с токеном из бд (bcrypt хэшом), если успешно, генерирует новую пару токенов для GUID и IP с которого отправлен запрос, и отправляет в формате json с полями access_token и refresh_token
  
3)/deleteToken
  Принимает json с полями: _id (GUID), access_token, refresh_token, делает проверку refresh токена пользователя сравнивая с токеном из бд (bcrypt хэшом), если успешно, удаляет запись с refresh токеном из бд по GUID и IP, отправляет статус 200, если проверка refresh токена не пройдена, то статус 401 
  
4)/deleteTokenы
  Принимает json с полями: _id (GUID), access_token, refresh_token, делает проверку refresh токена пользователя сравнивая с токеном из бд (bcrypt хэшом), если успешно, удаляет все записи с refresh токеном из бд по GUID и всех IP, отправляет статус 200, если проверка refresh токена не пройдена, то статус 401
