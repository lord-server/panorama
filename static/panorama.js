let map = L.map('map', {
    attributionControl: false,
    crs: L.CRS.Simple,
}).setView([0, 0], 0);
L.tileLayer('/tiles/{z}/{x}/{y}.png', {
    maxZoom: 0,
    minZoom: -7,
    tileSize: 256,
    noWrap: true,
}).addTo(map);
