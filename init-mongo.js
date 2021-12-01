db.createUser({
    user: "admin",
    pwd: "ASjkjasd13123",
    roles: [
        { role: "userAdminAnyDatabase", db: "admin" },
        { role: "readWriteAnyDatabase", db: "admin" },
        { role: "dbAdminAnyDatabase",   db: "admin" }
     ],
})
