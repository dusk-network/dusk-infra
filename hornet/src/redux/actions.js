import { ADD_NODE_UPDATE, ADD_REGION, CONNECTED, CONNECTING, CONNECTION_ERROR, DISCONNECTED, UPDATE_CPU_READ, UPDATE_DISK_READ, UPDATE_LAST_BLOCK_INFO, UPDATE_LOG_READ, UPDATE_MEM_READ, UPDATE_NET_READ, UPDATE_TIME_READ, UPDATE_WARN_LIST } from "./action-types";

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

export const updateBlockTimeRead = (value, timestamp) => ({
  type: UPDATE_TIME_READ,
  value,
  timestamp
})

export const updateWarningList = (warnings, timestamp) => ({
  type: UPDATE_WARN_LIST,
  warnings,
  timestamp
})

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
    console.log(data)
    const payload = JSON.parse(JSON.parse(data)); // Todo: fix wrong json encoding from server
    console.log(payload);

    const { metric, value, data:packet, timestamp } = payload;
    switch (metric) {
      case "cpu":
        dispatch(updateCPURead(parseFloat(value), getTime(timestamp)));
        break;
      case "mem":
        dispatch(updateMemoryRead(parseInt(value), getTime(timestamp)));
        break;
      case "latency":
        dispatch(updateNetRead(parseInt(value), getTime(timestamp)));
        break;
      case "disk":
        dispatch(updateDiskRead(parseInt(value), getTime(timestamp)));
        break;
      case "log":
        const { code } = packet
        if(code && code === "round"){
          const {round, blockHash, blockTime} = packet

          const block = { height: round, hash: blockHash, timestamp: timestamp }
          dispatch(updateLastBlockInfo(block));
          dispatch(updateBlockTimeRead(blockTime, getTime(timestamp)));
          break;
        };

        if(code && code === "warn"){
          const {time} = packet
          dispatch(updateWarningList(packet.data, getTime(time)))
          break;
        }

      case "tail":
        dispatch(updateLogRead(value, getTime(timestamp)));
        break;

      default:
        console.log("INVALID METRIC SENT");
    }
    // dispatch(addNodeUpdate(payload));
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
