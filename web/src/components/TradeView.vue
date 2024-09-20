<template>
  <div class="p-6 bg-gray-100 min-h-screen">
    <h2 class="text-2xl font-bold text-center text-gray-700 mb-6">Tokenize-Trade (BTC-ETH)</h2>
    <div class="grid grid-cols-2 gap-6 max-w-4xl mx-auto">
      <!-- Bid List -->
      <div class="bg-green-100 p-4 rounded-lg shadow-lg">
        <h3 class="text-xl font-semibold text-green-700 mb-4">Bids Summary: {{ bidSum }}</h3>
        <ul class="space-y-2">
          <li class="flex justify-around bg-white p-2 rounded shadow-sm text-green-900">
            <span>Size</span>
            <span>Price</span>
          </li>
          <li
              v-for="(order, index) in bidList"
              :key="index"
              class="flex justify-around bg-white p-2 rounded shadow-sm text-green-900"
          >
            <span>{{ order.size }}</span>
            <span>{{ order.price }}</span>
          </li>
        </ul>
      </div>

      <!-- Ask List -->
      <div class="bg-red-100 p-4 rounded-lg shadow-lg">
        <h3 class="text-xl font-semibold text-red-700 mb-4">Asks Summary: {{ askSum }}</h3>
        <ul class="space-y-2">
          <li class="flex justify-around bg-white p-2 rounded shadow-sm text-red-900">
            <span>Price</span>
            <span>Size</span>
          </li>
          <li
              v-for="(order, index) in askList"
              :key="index"
              class="flex justify-around bg-white p-2 rounded shadow-sm text-red-900"
          >
            <span>{{ order.price }}</span>
            <span>{{ order.size }}</span>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      bidSum: null,
      askSum: null,
      bidList: [],
      askList: [],
      socket: null,
      pingInterval: null,
    };
  },
  methods: {
    connectWebSocket() {
      const wsUrl = '/ws/symbol-depth';
      this.socket = new WebSocket(wsUrl);

      // 當接收到數據時，更新訂單簿
      this.socket.onmessage = (event) => {
        const data = JSON.parse(event.data);

        // 處理 bids
        let bidList = data.bids.map(bid => ({
          price: parseFloat(bid.price).toFixed(5),
          size: parseFloat(bid.size).toFixed(5),
        }));

        // 處理 asks
        let askList = data.asks.map(ask => ({
          price: parseFloat(ask.price).toFixed(5),
          size: parseFloat(ask.size).toFixed(5),
        }));

        // 更新數據到 Vue 的狀態
        this.bidSum = parseFloat(data.bids_sum).toFixed(5);
        this.askSum = parseFloat(data.asks_sum).toFixed(5);
        this.bidList = bidList;
        this.askList = askList;
      };

      this.socket.onerror = (error) => {
        console.error('WebSocket Error:', error);
      };

      this.pingInterval = setInterval(() => {
        if (this.socket.readyState === WebSocket.OPEN) {
          this.socket.send('ping');
        }
      }, 10000);
    },
  },
  mounted() {
    this.connectWebSocket();
  },
  beforeUnmount() {
    if (this.socket) {
      this.socket.close();
    }
    if (this.pingInterval) {
      clearInterval(this.pingInterval);
    }
  },
};
</script>
