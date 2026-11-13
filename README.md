# Go HTTP Server (CRUD API)

Go-এর স্ট্যান্ডার্ড লাইব্রেরি (`net/http`) এবং Go 1.22+ এর নতুন `http.ServeMux` রাউটিং ফিচার ব্যবহার করে তৈরি একটি সিম্পল RESTful HTTP সার্ভার।

---

## 📁 প্রজেক্ট স্ট্রাকচার

```text
go_server/
├── explore/
│   └── explore.go      # Go স্লাইস (Slice) থেকে এলিমেন্ট ডিলিট করার প্র্যাকটিস ফাইল
├── .gitignore          # গিট ইগনোর ফাইল
├── GEMINI.md           # প্রজেক্ট গাইডলাইন ও নেমিং কনভেনশন
├── go.mod              # Go মডিউল ফাইল
├── main.go             # সার্ভার সেটআপ, রাউটিং এবং সব API হ্যান্ডলার
└── README.md           # প্রজেক্ট ডকুমেন্টেশন
```

---

## ⚙️ প্রজেক্টটি কীভাবে কাজ করে?

### ১. সার্ভার সেটআপ (`main.go`)
- `http.NewServeMux()` ব্যবহার করে একটি রিকোয়েস্ট রাউটার তৈরি করা হয়েছে।
- পোর্ট `:5000`-এ সার্ভার লিসেন ও সার্ভ করে (`http.ListenAndServe(":5000", mux)`)।

### ২. রাউটিং মেকানিজম (Go 1.22+ Routing)
Go 1.22 থেকে সরাসরি রাউটের ভেতরেই HTTP Method এবং ডায়নামিক পাথ প্যারামিটার হ্যান্ডেল করা যায়:
- **মেথড ম্যাচিং:** যেমন `POST /create-user`, `GET /user/{id}`, `DELETE /delete-user` ইত্যাদি।
- **পাথ প্যারামিটার:** `r.PathValue("id")` ব্যবহার করে ইউআরএল থেকে ভ্যালু সংগ্রহ করা হয়।

### ৩. ইন-মেমোরি ডেটা স্টোরেজ
ডেটাবেজের পরিবর্তে একটি গ্লোবাল স্লাইস `users` ব্যবহার করে ডেটা স্টোর করা হয়েছে:
```go
var users = []User{
    {ID: 1, Name: "Rahim", Email: "rahim@example.com"},
    {ID: 2, Name: "Karim", Email: "karim@example.com"},
    {ID: 3, Name: "Kobirul", Email: "kobirul@example.com"},
}
```

---

## 🚀 API এন্ডপয়েন্টসমূহ

| Method | Endpoint | Description | Status Code |
| :--- | :--- | :--- | :--- |
| `GET` | `/` | রুট হ্যান্ডলার (স্বাগতম বার্তা) | `200 OK` |
| `GET` | `/health` | হেলথ চেক এন্ডপয়েন্ট | `200 OK` |
| `GET` | `/api/users` | সকল ইউজারের তালিকা রিটার্ন করে (JSON) | `200 OK` |
| `GET` | `/user/{id}` | নির্দিষ্ট ID-র ইউজারের তথ্য রিটার্ন করে | `200 OK` / `404 Not Found` |
| `POST` | `/create-user` | নতুন ইউজার তৈরি করার এন্ডপয়েন্ট | `200 OK` |
| `POST` | `/create-post` | নতুন পোস্ট তৈরি করার এন্ডপয়েন্ট | `200 OK` |
| `PUT` | `/update-user` | ইউজার সম্পূর্ণ আপডেট/রিপ্লেস করার জন্য | `200 OK` |
| `PUT` | `/update-post` | পোস্ট আপডেট/রিপ্লেস করার জন্য | `200 OK` |
| `PATCH` | `/update-user` | ইউজারের আংশিক তথ্য আপডেট করার জন্য | `200 OK` |
| `DELETE` | `/delete-user/{id}` | নির্দিষ্ট ID-র ইউজার ডিলিট করার এন্ডপয়েন্ট | `204 No Content` / `404 Not Found` |

---

## 🛠️ সার্ভার রান করার নিয়ম

টার্মিনালে প্রজেক্ট ডিরেক্টরিতে গিয়ে রান করুন:

```bash
go run .
```
সার্ভার সফলভাবে চালু হলে কনসোলে দেখাবে:
```text
Server is running on port 5000...
```

---

## 🧪 API টেস্ট করার কিছু উদাহরণ (cURL)

#### ১. সব ইউজার দেখতে:
```bash
curl http://localhost:5000/api/users
```

#### ২. নির্দিষ্ট ইউজার দেখতে (ID: 1):
```bash
curl http://localhost:5000/user/1
```

#### ৩. নতুন ইউজার তৈরির রিকোয়েস্ট (POST):
```bash
curl -X POST http://localhost:5000/create-user
```

#### ৪. হেলথ চেক:
```bash
curl http://localhost:5000/health
```

---

## 🔍 `explore/explore.go` ফাইলে কী হচ্ছে?
Go স্লাইস থেকে নির্দিষ্ট ইনডেক্সের উপাদান ডিলিট করার টেকনিক প্র্যাকটিস করা হয়েছে:
```go
// idx নম্বর উপাদান বাদ দিয়ে বাকি দুটি অংশ যুক্ত করা হয়:
users = append(users[:idx], users[idx+1:]...)
```

