async function auth() {
    const login = document.getElementById('login').value
    const password = document.getElementById('password').value
    const link = document.getElementById('link').value

    const user = { login, password, link }

    const response = await fetch('http://localhost:9000/auth', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json;charset=utf-8'
        },
        body: JSON.stringify(user)
    })

    const result = await response.json()
    alert(result.message)
}
