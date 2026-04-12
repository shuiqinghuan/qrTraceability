package main

import (
    "fmt"
    "log"

    "golang.org/x/crypto/bcrypt"
)

func main() {
    password := "123456"
    
    // 生成新哈希
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("密码: %s\n", password)
    fmt.Printf("新哈希: %s\n", string(hash))
    
    // 验证新哈希
    err = bcrypt.CompareHashAndPassword(hash, []byte(password))
    if err != nil {
        fmt.Println("验证失败!")
    } else {
        fmt.Println("验证成功!")
    }
}
