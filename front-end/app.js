import ForceGraph3D from '3d-force-graph';

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
      .backgroundColor('#262b36')
      .nodeAutoColorBy('type')
      .nodeVal('count')
      .linkDirectionalParticleSpeed(d => d.weight);

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
    
    let { nodes, links } = Graph.graphData();

    Graph.graphData({
      nodes: [...nodes, ...parsed.graph.nodes],
      links: [...links, ...parsed.graph.edges]
    });
  
  }
});
