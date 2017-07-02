<template>
  <div class="hello">
    <h1>Products</h1>
    <br>
    <br>
    <ul>
      <li v-for="p in products" :key="p.id">
        <button @click="deleteProduct(p.id)">X</button>
        <router-link :to="{ name: 'Product', params: { id: p.id }}">
          {{ p.name }}
        </router-link>&nbsp;
        {{ p.price }} €
      </li>
    </ul>
    <br>
    <hr>
    <br>
    <form>
      <input v-model="productName" placeholder="Product name">
      <br>
      <br>
      <input type="number" number v-model="productPrice" placeholder="Price">
      <br>
      <br>
      <input type="submit" value="Ajouter" @click="createProduct()">
    </form>
  </div>
</template>

<script>
import axios from 'axios';

const api = '/api'; // 'http://localhost:5678/api';

export default {
  name: 'hello',
  data: () => ({
    productName: null,
    productPrice: 0.0,
    products: [],
  }),

  methods: {
    async createProduct() {
      await axios.post(`${api}/products`, {
        name: this.productName,
        price: Number(this.productPrice),
      });

      // refresh the data
      this.retrieveProducts();
    },

    async deleteProduct(id) {
      // delete the product
      await axios.delete(`${api}/products/${id}`);

      // refresh the data
      const response = await axios.get(`${api}/products`);
      this.products = response.data;
    },

    async retrieveProducts() {
      const response = await axios.get(`${api}/products`);
      this.products = response.data.sort((a, b) => a.id - b.id);
    },
  },

  async created() {
    this.retrieveProducts();
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1, h2 {
  font-weight: normal;
}

ul {
  list-style-type: none;
  padding: 0;
  padding-left: 40%;
}

li {
  /*display: inline-block;*/
  margin: 0 10px;
  text-align: left;
}

a {
  color: #42b983;
}
</style>
