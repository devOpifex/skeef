import ForceGraph3D from '3d-force-graph';
import SpriteText from 'three-spritetext';

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

  let graphElem = document.getElementById("graph");
  let Graph;
  if(graphElem != null){
    const initData = {
      nodes: [],
      links: []
    };

    Graph = ForceGraph3D()(graphElem)
      .onNodeHover(node => graphElem.style.cursor = node ? 'pointer' : null)
      .graphData(initData)
      .enableNodeDrag(false)
      .backgroundColor('#262b36')
      .nodeAutoColorBy('type')
      .nodeVal('value')
      .nodeThreeObject(node => {
        const sprite = new SpriteText(node.id);
        sprite.material.depthWrite = false; // make sprite background transparent
        sprite.color = node.color;
        sprite.textHeight = node.value;
        return sprite;
      })
      .linkDirectionalParticles("value")
      .linkDirectionalParticleSpeed(d => d.value * 0.001);

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

  let firstrun = true;
  let ntweets = document.getElementById("ntweets");

  window.socket.onmessage = (data) => {
    let parsed = JSON.parse(data.data);
    console.log(parsed);
    
    if(parsed.tweetsCount == 0){
      return ;
    }

    ntweets.innerText = parsed.tweetsCount;

    let { nodes, links } = Graph.graphData();

    if (firstrun == true){
      Graph.graphData({
        nodes: [...nodes, ...parsed.graph.nodes],
        links: [...links, ...parsed.graph.edges]
      });
      firstrun = false;
      return ;
    }

    // only nodes to add
    let nodesAdditions = parsed.graph.nodes.filter(function(n){
      return n.action == "add"
    })

    if(parsed.graph.nodes != null){
      let nodeUpdates = parsed.graph.nodes.filter(function(n){
        return n.action == "update"
      });
  
      for(let i = 0; i < nodeUpdates.length; i++){
        for(let j = 0; j < nodes.length; j++){
          if(nodes[j].id == nodeUpdates[i].name){
            nodes[j].value = nodeUpdates[i].value
          }
        }
      }
    }

    let edgesAdditions = parsed.graph.edges.filter(function(e){
      return e.action == "add"
    });

    if(parsed.graph.nodes != null){
      let edgeUpdates = parsed.graph.nodes.filter(function(e){
        return e.action == "update"
      });
  
      for(let i = 0; i < edgeUpdates.length; i++){
        for(let j = 0; j < links.length; j++){
          if(links[j].id == edgeUpdates[i].name){
            links[j].value = edgeUpdates[i].value
          }
        }
      }
    }

    // remove
    if(parsed.removeNodes.length > 0){
      for(let node of parsed.removeNodes){
        links = links.filter(l => l.source !== node && l.target !== node);
        nodes = nodes.filter(n => n.id != node);
      }
    }

    Graph.graphData({
      nodes: [...nodes, ...nodesAdditions],
      links: [...links, ...edgesAdditions]
    });
  
  }
});
