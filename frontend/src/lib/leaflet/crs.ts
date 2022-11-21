import * as L from 'leaflet';

export const isometric = L.extend(L.CRS.Simple, {
    projection: {
        project: (latlng: L.LatLng) => {
            let x = ((latlng.lng - latlng.lat) * 16) / 2;
            let y = ((latlng.lng + latlng.lat) * 16) / 4;
            return new L.Point(x, y);
        },

        unproject: (point: L.Point) => {
            point.y *= 2;
            let lat = (point.y - point.x) / 16;
            let lng = (point.y + point.x) / 16;
            return new L.LatLng(lat, lng);
        }
    },
    transformation: new L.Transformation(1, 0, 1, 0)
});

export const flat = L.extend(L.CRS.Simple, {
    projection: {
        project: (latlng: L.LatLng) => {
            let x = latlng.lat * 16;
            let y = latlng.lng * 16;
            return new L.Point(x, y);
        },

        unproject: (point: L.Point) => {
            let lat = point.x / 16;
            let lng = point.y / 16;
            return new L.LatLng(lat, lng);
        }
    },
    transformation: new L.Transformation(1, 0, 1, 0)
});
