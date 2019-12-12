<template>
  <div class="grid">
    <div class="grid__item theme-1">
      <div class="grid__item-content">
        <nav class="grid__item-nav">
          <a
            href="https://tympanus.net/Development/GridMenu/"
            class="grid__item-link"
          >&larr; Previous Demo</a>
          <a href="https://tympanus.net/codrops/?p=34430" class="grid__item-link">Article</a>
          <!--a href="" class="grid__item-link">GitHub</a-->
        </nav>
        <h2 class="grid__item-title grid__item-title--small">Distortion Hover Effect</h2>
      </div>
    </div>
    <div class="hello">
      <div class="grid__item grid__item--bg theme-2">
        <div class="grid__item-img"></div>
        <div class="grid__item-content">
          <span class="grid__item-meta">California</span>
          <h2 class="grid__item-title">Jumping Around</h2>
          <h3 class="grid__item-subtitle">
            <span>California's last empty jump spots</span>
            <a class="grid__item-link" href="#">Discover more</a>
          </h3>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import hoverEffect from "hover-effect";

export default {
  name: "HelloWorld",
  props: {
    msg: String
  },
  methods: {},
  mounted() {
    console.log("ttttttttttt");
    var myAnimation = new hoverEffect({
      parent: document.querySelector(".grid__item-img"),
      intensity: 0.3,
      image1: "/img/Img22.jpg",
      image2: "/img/Img21.jpg",
      displacementImage: "/img/displacement/1.jpg"
    });
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}

@import "../css/base.css";
</style>
