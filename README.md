# booking
Prove of concept ticket booking system implement by grpc

client state
0 start
    - not open -> 0
    - open -> 1
1 wait room
    - ready -> 2
2 login
    - success -> 3
3 booking
    - select reserved seating -> 3.1
    - select general admission ->3.2
    - no free seat -> 5
    3.1 booking reserve - select seat
        - success -> 4
        - fail -> 3
    3.2 booking GA - number of seat
        - success -> 4
        - fail -> 3
4 payment
    - sucess -> 5
    - fail -> 4
    - not pay -> 5
5 end

reserved seat state
1 free
    - book -> 2
2 book
    - pay -> 3
    - timeout ->1
3 paid


