console.log("Hello, world");

document.addEventListener("DOMContentLoaded",function(){
  let locationsInput = document.getElementById("locationsStream");
  let locationsWarning = document.getElementById("locationsWarning");
  locationsInput.addEventListener("input", (event) => {
    locationsWarning.innerHTML = "This dramatically reduces the number of tweets streamed, see the documentation."
  });
});