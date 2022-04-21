L.Control.CursorPosition = L.Control.extend({
    options: {
        position: 'topleft',
    },
    onAdd: map => {
        return L.DomUtil.create('div', 'cursor-position leaflet-control leaflet-bar');
    }
});

L.CRS.Isometric = L.extend(L.CRS.Simple, {
    projection: {
        project: latlng => {
            let x = (latlng.lng - latlng.lat) * 16 / 2;
            let y = (latlng.lng + latlng.lat) * 16 / 4;
            return new L.Point(x, y);
        },

        unproject: point => {
            point.y *= 2;
            let lat = (point.y - point.x) / 16;
            let lng = (point.y + point.x) / 16;
            return new L.LatLng(lat, lng);
        }
    },
    transformation: new L.Transformation(1, 0, 1, 0),
});

class Map {
    map;

    constructor(zoomLevels) {
        this.map = L.map('map', {
            attributionControl: false,
            crs: L.CRS.Isometric,
        }).setView([735, -818], 0);

        let position = new L.Control.CursorPosition();
        position.addTo(this.map);

        this.map.on('click mouseover mousemove', e => {
            let x = e.latlng.lat;
            let z = e.latlng.lng;
            position.getContainer().innerHTML = `X=${x.toFixed()} Y=0 Z=${z.toFixed()}`;
        });

        L.tileLayer('/tiles/{z}/{x}/{y}.png', {
            maxZoom: 0,
            minZoom: -zoomLevels,
            tileSize: 256,
            noWrap: true,
        }).addTo(this.map);
    }
}

async function main() {
    let meta = await fetch('/metadata.json').then(response => response.json());
    document.title = meta['title'];

    let map = new Map(parseInt(meta['zoomLevels'], 10));
}

main();

