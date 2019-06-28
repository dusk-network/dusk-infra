import {
  ADD_NODE_UPDATE,
  UPDATE_CPU_READ,
  UPDATE_MEM_READ,
  UPDATE_NET_READ,
  UPDATE_LOG_READ,
  UPDATE_DISK_READ,
  ADD_REGION,
  UPDATE_LAST_BLOCK_INFO,
  CONNECTING,
  CONNECTED,
  CONNECTION_ERROR,
  DISCONNECTED
} from "./action-types";

export const addNodeUpdate = payload => ({
  type: ADD_NODE_UPDATE,
  payload
});

export const updateCPURead = (value, timestamp) => ({
  type: UPDATE_CPU_READ,
  value,
  timestamp
});

export const updateNetRead = (value, timestamp) => ({
  type: UPDATE_NET_READ,
  value,
  timestamp
});

export const updateDiskRead = (value, timestamp) => ({
  type: UPDATE_DISK_READ,
  value,
  timestamp
});

export const updateLogRead = (value, timestamp) => ({
  type: UPDATE_LOG_READ,
  value,
  timestamp
});

export const updateMemoryRead = (value, timestamp) => ({
  type: UPDATE_MEM_READ,
  value,
  timestamp
});

export const addRegion = payload => ({
  type: ADD_REGION,
  payload
});

export const updateLastBlockInfo = payload => ({
  type: UPDATE_LAST_BLOCK_INFO,
  payload
});

export const connecting = payload => ({
  type: CONNECTING,
  payload
});

export const connected = payload => ({
  type: CONNECTED,
  payload
});

export const connectionError = payload => ({
  type: CONNECTION_ERROR,
  payload
});

export const disconnected = payload => ({
  type: DISCONNECTED,
  payload
});

export const listenForUpdates = socket => dispatch => {
  dispatch(connecting());
  let ws = new WebSocket("ws://localhost:8080/stats");

  ws.onopen = () => dispatch(connected());
  ws.onerror = () => dispatch(connectionError());
  ws.onclose = () => dispatch(disconnected());
  ws.onmessage = ({ data }) => {
    const payload = JSON.parse(JSON.parse(data)); // Todo: fix wrong json encoding from server
    console.log(data);
    console.log(payload);

    const { metric, value, timestamp } = payload;
    switch (metric) {
      case "cpu":
        dispatch(updateCPURead(parseFloat(value), getTime(timestamp)));
        break;
      case "mem":
        dispatch(updateMemoryRead(parseInt(value), getTime(timestamp)));
        break;
      case "net":
        dispatch(updateNetRead(parseInt(value), getTime(timestamp)));
        break;
      case "dsk":
        dispatch(updateDiskRead(parseInt(value), getTime(timestamp)));
        break;
      case "log":
        dispatch(updateLogRead(value, getTime(timestamp)));
        break;

      default:
        console.log("INVALID METRIC SENT");
    }
    // dispatch(addNodeUpdate(payload));
    // dispatch(updateLastBlockInfo(payload));
    // dispatch(addRegion(payload));
  };

  if (!("ws" in socket)) {
    socket.ws = ws;
  }
};

const getTime = timestamp => timestamp;
// let g = new Date(Date.parse(timestamp));
// let hours = g.getHours();
// let seconds = g.getSeconds();
// let minutes = g.getMinutes();
// return (
//   hours +
//   ":" +
//   minutes.toString().padStart(2, "0") +
//   ":" +
//   seconds.toString().padStart(2, "0")
// );
