document.getElementById("form")
    .addEventListener("submit", async e => {
        e.preventDefault()

        const data = new FormData(e.target)
        const value = Object.fromEntries(data.entries())

        const msg = await fetch(`/msg/${value.messageId}`,
            {
                method: "GET",
            })
            .then(resp => resp.json())
            .then(body => JSON.stringify(body, null, 4))

        document.getElementById("msg-view").innerText = msg
    })