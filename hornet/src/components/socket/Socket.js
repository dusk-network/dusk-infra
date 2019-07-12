import { connect } from "react-redux";
import {
  connectTo,
  connected,
  connectionError,
  disconnected,
  updateLastBlockInfo,
  updateBlockTimeRead,
  updateWarningList,
  updateLogRead,
  updateCPURead,
  updateMemoryRead,
  updateDiskRead,
  updateNetRead,
} from "../../redux/actions";
import { useEffect, useRef } from "react";
import { isThisWeek } from "date-fns";
//
const updateMetrics = {
  cpu: updateCPURead,
  mem: updateMemoryRead,
  latency: updateNetRead,
  disk: updateDiskRead,
};

class DuskSocket {
  constructor(dispatch) {
    this.dispatch = dispatch;
    this.ws = null;
  }

  open(hostname, port) {
    const { dispatch } = this;
    this.ws = new WebSocket(`ws:/${hostname}:${port}/stats`);

    ["open", "close", "error", "message"].forEach(type =>
      this.ws.addEventListener(type, this)
    );
  }

  handleEvent(event) {
    const listener = this[`on${event.type}`];
    listener && listener.call(this, event);
  }

  onmessage({ data }) {
    const { dispatch } = this;

    const payload = JSON.parse(data);
    const { metric, value, data: packet, timestamp } = payload;

    switch (metric) {
      case "cpu":
      case "mem":
      case "latency":
      case "disk":
        dispatch(updateMetrics[metric](+value, timestamp));
        break;
      case "log":
        const { code, level } = packet;
        if (code && code === "round") {
          const { round, blockHash, blockTime } = packet;
          const block = {
            height: round,
            hash: blockHash,
            timestamp,
          };
          dispatch(updateLastBlockInfo(block));
          dispatch(updateBlockTimeRead(blockTime, timestamp));
          break;
        }
        if (level) {
          const { time = timestamp } = packet;
          dispatch(updateWarningList(packet, time));
          break;
        }
      case "tail":
        dispatch(updateLogRead(value, timestamp));
        break;
      default:
        console.log("INVALID METRIC SENT");
    }
  }

  onopen() {
    this.dispatch(connected());
  }

  onclose() {
    this.dispatch(disconnected());
  }

  onerror() {
    this.dispatch(connectionError());
  }

  close() {
    this.ws.close();
  }
}

const Socket = ({ connectTo, status, hostname, port, dispatch }) => {
  const ws = useRef();

  useEffect(() => {
    const [hostname, port] = (
      process.env.REACT_APP_HOST_WS || window.location.host
    ).split(":");

    connectTo(hostname, port);

    return () => {
      console.log("cleanup");
    };
  }, []);

  useEffect(() => {
    switch (status) {
      case "connecting":
        ws.current = new DuskSocket(dispatch);
        ws.current.open(hostname, port);
    }
  }, [status]);

  return null;
};

const mapStateToProps = ({ status, hostname, port }) => ({
  status,
  hostname,
  port,
});

const mapDispatchToProps = dispatch => ({
  connectTo: (hostname, port) => dispatch(connectTo(hostname, port)),
  dispatch,
});

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Socket);
