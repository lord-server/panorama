let map = L.map('map', {
    crs: L.CRS.Simple,
}).setView([0, 0], 0);
L.tileLayer('/tiles/{x}/{y}.png', {
    attribution: '<a href="https://github.com/lord-server/panorama">panorama</a>',
    maxZoom: 0,
    tileSize: 256,
    noWrap: true,
}).addTo(map);
