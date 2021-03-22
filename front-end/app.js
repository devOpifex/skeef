document.addEventListener("DOMContentLoaded",function(){

  // warning on admin inputs
  let locationsInput = document.getElementById("locationsStream");
  let locationsWarning = document.getElementById("locationsWarning");
  locationsInput.addEventListener("input", (event) => {
    if(event.data == null){
      locationsWarning.innerHTML = "";
      return ;
    }
    locationsWarning.innerHTML = "Few tweets are geo-tagged; this dramatically reduces the number of tweets streamed."
  });

  // websocket
  window.socket.onopen = () => {
    socket.send("Client Connected")
  }
  
  window.socket.onclose = () => {
    console.log("Websocket closed")
  }
  
  window.socket.onerror = (error) => {
    console.log(error);
  }
  
  window.socket.onmessage = (data) => {
    let parsed = JSON.parse(data.data);
    console.log(parsed);
  }
});
