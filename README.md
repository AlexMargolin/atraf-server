# TODO

1. ~~remove mysql from handlers~~
2. ~~give NewServices a better name~~
3. Validate CORS (origin)
4. Validate CSRF
5. ~~Better error handling~~
6. ~~Remove Pkg imports from main. should only contain internals~~
7. Storage services are passed as references
8. System Load test
9. Generate strong jwt secret
10. change user_id to account_id
11. define request tags(?uuid and so on i guess)
12. comments pagination
13. security headers
14. owasp
15. ~~remove config crap and use simple os.env with checker~~
16. validate jwt alg
17. session context - change accountid to userid
18. better comments
19. test user inputs
20. ~~refactor pagination = replace offset with cursor~~
21. server logger
22. enable server tls
23. map json response values to 2 letters
24. create updated at trigger
25. create indexes on created_at
26. change {Domain}Fields to Data
27. todo remove mail from public users
28. mess test pagination base64 json
29. invoke jwt
30. create jwt aud for password reset
31. move login access token into service
32. add updated at trigger
33. foreign keys (account + account_reset)
34. cron job to delete requested password resets
35. function comments
36. prevent time attacks @ login / register and any account related stuff
37. email lower case
38. try using different JWTs across different endpoints

# Features

1. [Posts] Attachments
2. [Auth] Refresh Token
3. [Reactions]