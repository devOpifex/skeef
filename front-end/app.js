var createGraph = require('ngraph.graph');
var pixel = require('ngraph.pixel');

document.addEventListener("DOMContentLoaded",function(){

  // warning on admin inputs
  let locationsInput = document.getElementById("locationsStream");
  
  if(locationsInput != null){
    let locationsWarning = document.getElementById("locationsWarning");
    locationsInput.addEventListener("input", (event) => {
      if(event.data == null){
        locationsWarning.innerHTML = "";
        return ;
      }
      locationsWarning.innerHTML = "Few tweets are geo-tagged; this dramatically reduces the number of tweets streamed."
    });
  }

  var g = createGraph();
  g.addNode(1, {type: 'hidden'});
  g.addNode(2, {type: 'hidden'});

  if(document.getElementById("graph") != null){
    var renderer = pixel(g, {
      is3d: true,
      container: document.getElementById("graph"),
      physics: {
        springLength : 80,
        springCoeff : 0.0002,
        gravity: -1.2,
        theta : 0.8,
        dragCoeff : 0.02
      },
      node: function (n) {
        if(n.data.type == 'user'){
          return {
            color: 0xFFC6FF,
            size: n.data.count * 10
          }
        } else if(n.data.type == 'hashtag'){
          return {
            color: 0x9BF6FF,
            size: n.data.count * 10
          }
        } else if (n.data.type == 'hidden') {
          return {
            color: 0x000000,
            size: 0
          }
        } else {
          return {
            color: 0xFFFFFF,
            size: 15
          }
        }
      },
      link: function(l){
        return {
          fromColor: 0xCCCCCC,
          toColor: 0xCCCCCC
        }
      }
    });

    renderer.clearColor("0x262b36");
  }

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

    g.beginUpdate();
    if(parsed.graph.nodes){
      for(let i = 0; i < parsed.graph.nodes.length; i++)  {

        if(g.hasNode(parsed.graph.nodes[i].name)){
          continue;
        }

        g.addNode(
          parsed.graph.nodes[i].name,
          {type: parsed.graph.nodes[i].type, count: parsed.graph.nodes[i].count}
        );
      }
    }

    if(parsed.graph.edges){
      for(let i = 0; i < parsed.graph.edges.length; i++){
        if(g.hasLink(parsed.graph.edges[i].source, parsed.graph.edges[i].target)){
          continue; 
        }

        g.addLink(parsed.graph.edges[i].source, parsed.graph.edges[i].target);
      }
    }
    g.endUpdate();
  }
});
