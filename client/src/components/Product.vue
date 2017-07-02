<template>
<div>
    <h1>{{ name }}</h1>
    <form>
      <input type="text" v-model="name">
      <br><br>
      <input number type="number" v-model="price"> €
      <br><br>
      <input type="submit" value="Edit" @click="updateProduct()">
      <br><br>
    </form>
    <router-link to="/">
      Retour
    </router-link>
</div>
</template>

<script>
import axios from 'axios';

const api = '/api'; // 'http://localhost:5678/api';

export default {
  name: 'product',
  data: () => ({
    name: null,
    id: null,
    price: null,
  }),

  async created() {
    const result = await axios.get(`${api}/products/${this.$route.params.id}`);

    this.id = result.data.id;
    this.name = result.data.name;
    this.price = result.data.price;
  },

  methods: {
    async updateProduct() {
      await axios.put(`${api}/products/${this.id}`, {
        name: this.name,
        id: this.id,
        price: Number(this.price),
      });
    },
  },
};
</script>

<style scoped>
a {
  color: #42b983;
}
</style>
