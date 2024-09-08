const evtSource = new EventSource("//localhost:3012/api/events", {
    withCredentials: true,
});
evtSource.onmessage = (event) => {
    if(event.data === '__RELOAD__') {
        window.location.reload()
    }
};
evtSource.onerror = (err) => {
    console.error("EventSource failed:", err);
};
