package main

import (
    "os"
    "fmt"
    "strconv"
    "log"
    "errors"
    "net/http"
    "encoding/json"
    "bytes"
    "database/sql"
    _ "github.com/lib/pq"
)

var db *sql.DB

type Response struct {
    Message string `json:"message"`
}

type User struct {
    Name string `json:"name"`
    Email string `json:"email"`
}

type Data struct {
    Id int `json:"id"`
    Name string `json:"name"`
    Email string `json:"email"`
    Created_at string `json:"created_at"`
    Updated_at string `json:"updated_at"`
}

const (
    DB_USER = "postgres"
    DB_PASSWORD = "postgres"
    DB_NAME = "postgres"
    TABLE_NAME = "kadaidb"
)

func struct2jsonstr (d interface{}) (string, error) {
    jsonBytes, err := json.Marshal(d)
    if err != nil {
        fmt.Println("JSON Marshal error:", err)
        return "", err
    }
    out := new(bytes.Buffer)
    json.Indent(out, jsonBytes, "", "    ")
    return out.String(), nil
}

func getUser (w http.ResponseWriter, r *http.Request) (User, error) {
    bufbody := new(bytes.Buffer)
    bufbody.ReadFrom(r.Body)
    body := bufbody.String()
    bodyBytes := ([]byte)(body)
    user := new(User)
    err := json.Unmarshal(bodyBytes, user)
    if err != nil && (len(user.Name) == 0 || len(user.Email) == 0) {
        err = errors.New("Try like {\"name\":\"a\", \"email\":\"test@a\"}")
    }
    if err != nil {
        fmt.Fprintln(w, "Invalid User Data:", err)
        return *user, err
    } else {
        return *user, nil
    }
}

func main() {
    var err error
    port := os.Getenv("PORT")
    if port == "" {
        port = "5432"
    }
    dbinfo := fmt.Sprintf("postgres://%s:%s@postgres:%s/%s?sslmode=disable", DB_USER, DB_PASSWORD, port, DB_NAME)
    url := os.Getenv("DATABASE_URL")
    if url != "" {
        dbinfo = url
    }
    db, err = sql.Open("postgres", dbinfo)
    if err != nil {
        fmt.Println("sql.Open")
        panic(err)
    }
    defer db.Close()

    tablecreate := fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s (
        id SERIAL NOT NULL,
        name TEXT NOT NULL,
        email TEXT NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    )
    `, TABLE_NAME)
    _, err = db.Exec(tablecreate)
    if err != nil {
        panic(err)
    }

    http.HandleFunc("/users/", useridRes)
    http.HandleFunc("/users", userRes)
    http.HandleFunc("/", helloRes)

    log.Fatal(http.ListenAndServe(":80", nil))
}

func helloRes (w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    switch r.Method {
    case "GET":
        w.Header().Add("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        res := Response{Message:"Hello World!!"}
        jres, err := struct2jsonstr(res)
        if err != nil {
            http.Error(w, "Internal Error.", 500)
            return
        }
        fmt.Fprintln(w, jres)
    default:
        http.Error(w, "Invalid request method.", 405)
        return
    }
}

func userRes (w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        w.Header().Add("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)

        GET := fmt.Sprintf("SELECT * FROM %s ORDER BY id", TABLE_NAME)
        rows, err := db.Query(GET)
        if err != nil {
            http.Error(w, "Internal Error.", 500)
            return
        }
        for rows.Next() {
            var data Data
            err = rows.Scan(&data.Id, &data.Name, &data.Email, &data.Created_at, &data.Updated_at)
            if err != nil {
                return
            }
            jdata, err := struct2jsonstr(data)
            if err!= nil {
                http.Error(w, "Internal Error.", 500)
                return
            }
            fmt.Fprintln(w, jdata)
        }

    case "POST":
        w.Header().Add("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)

        user, err := getUser(w, r)
        if err != nil {
            http.Error(w, "Internal Error.", 500)
            return
        }
        var lastInsertId int
        INSERT := fmt.Sprintf("INSERT INTO %s (name,email,created_at) VALUES($1,$2,current_timestamp) returning id;", TABLE_NAME)
        err = db.QueryRow(INSERT, user.Name, user.Email).Scan(&lastInsertId)
        if err != nil {
            http.Error(w, "Internal Error.", 500)
            return
        }
        var data Data
        GET := fmt.Sprintf("SELECT * FROM %s WHERE id=$1;", TABLE_NAME)
        err = db.QueryRow(GET, lastInsertId).Scan(&data.Id, &data.Name, &data.Email, &data.Created_at, &data.Updated_at)
        if err != nil {
            http.Error(w, "Internal Error.", 500)
            return
        }
        jdata, err := struct2jsonstr(data)
        if err != nil {
            http.Error(w, "Internal Error.", 500)
            return
        }
        fmt.Fprintln(w, jdata)

    default:
        http.Error(w, "Invalid request method.", 405)
        return
    }
}

func useridRes (w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.URL.Path[len("/users/"):])
    if err != nil {
        http.Error(w, "Invalid ID.", 500)
        return
    }
    switch r.Method {
    case "GET":
        w.Header().Add("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)

        GET := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", TABLE_NAME)
        var data Data
        err = db.QueryRow(GET, id).Scan(&data.Id, &data.Name, &data.Email, &data.Created_at, &data.Updated_at)
        if err != nil {
            http.Error(w, "Internal Error.", 500)
            return
        }
        jdata, err := struct2jsonstr(data)
        if err != nil {
            http.Error(w, "Internal Error.", 500)
            return
        }
        fmt.Fprintln(w, jdata)

    case "PUT":
        w.Header().Add("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)

        user, err := getUser(w, r)
        if err != nil {
            http.Error(w, "Internal Error.", 500)
            return
        }
        UPDATE := fmt.Sprintf("UPDATE %s SET name=$1,email=$2,updated_at=current_timestamp WHERE id=$3 returning *", TABLE_NAME)
        var data Data
        err = db.QueryRow(UPDATE, user.Name, user.Email, id).Scan(&data.Id, &data.Name, &data.Email, &data.Created_at, &data.Updated_at)
        if err != nil {
            http.Error(w, "Internal Error.", 500)
            return
        }
        jdata, err := struct2jsonstr(data)
        if err != nil {
            http.Error(w, "Internal Error.", 500)
            return
        }
        fmt.Fprintln(w, jdata)

    case "DELETE":
        w.WriteHeader(http.StatusNoContent)

        DELETE := fmt.Sprintf("DELETE FROM %s WHERE id=$1", TABLE_NAME)
        _, err = db.Query(DELETE, id)
        if err != nil {
            http.Error(w, "Internal Error.", 500)
            return
        }

    default:
        http.Error(w, "Invalid request method.", 405)
        return
    }
}
