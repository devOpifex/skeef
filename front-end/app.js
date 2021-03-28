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
      clearAlpha: 0.0,
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
            color: 0xff9e00,
            size: n.data.count * 10
          }
        } else if(n.data.type == 'hashtag'){
          return {
            color: 0xb700ff,
            size: n.data.count * 10
          }
        } else if (n.data.type == 'hidden') {
          return ;
        }
      },
      link: function(l){
        return {
          fromColor: 0xCCCCCC,
          toColor: 0xCCCCCC
        }
      }
    });
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

  let ntweets = document.getElementById("ntweets")
  window.socket.onmessage = (data) => {
    let parsed = JSON.parse(data.data);
    console.log(parsed);
    
    if(parsed.tweetsCount != 0){
      ntweets.innerText = parsed.tweetsCount;
    }

    g.beginUpdate();
    if(parsed.graph.nodes){
      for(let i = 0; i < parsed.graph.nodes.length; i++)  {

        if(parsed.graph.nodes[i].action == "update")
          continue

        if(parsed.graph.nodes[i].action == "add"){
          g.addNode(
            parsed.graph.nodes[i].name,
            {type: parsed.graph.nodes[i].type, count: parsed.graph.nodes[i].count}
          );
        }

      }
    }

    if(parsed.graph.edges){
      for(let i = 0; i < parsed.graph.edges.length; i++){
        if(parsed.graph.edges[i].action == "update"){
          continue; 
        }

        if(parsed.graph.edges[i].action == "add"){
          g.addLink(parsed.graph.edges[i].source, parsed.graph.edges[i].target);
        }
        
      }
    }
    g.endUpdate();

    g.beginUpdate();
    if(parsed.graph.nodes){
      for(let i = 0; i < parsed.graph.nodes.length; i++)  {

        if(parsed.graph.nodes[i].action == "add")
          continue

        if(parsed.graph.nodes[i].action == "update"){
          let n = g.getNode(parsed.graph.nodes[i].name);

          // if undefined then we have an issue and we must add 
          // the node instead
          if(n == undefined) {
            g.addNode(
              parsed.graph.nodes[i].name,
              {type: parsed.graph.nodes[i].type, count: parsed.graph.nodes[i].count}
            );
            continue
          }

          n.count = parsed.graph.nodes[i].count;
        }

      }
    }
    g.endUpdate();
  }
});
