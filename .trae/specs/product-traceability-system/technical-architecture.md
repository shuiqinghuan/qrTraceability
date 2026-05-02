# 农产品溯源系统 - 技术架构文档（简化版）

## 1. 技术选型（极简方案）

**核心原则：**
- 最少依赖，最快部署
- 零后端，纯静态
- 单文件组件优先

**技术栈：**
```
Vue 3（CDN引入）
├── Vue 3 Core（CDN）
├── 纯CSS（无框架）
└── 静态JSON数据
```

**无依赖方案对比：**
| 方案 | 复杂度 | 部署难度 | 适用场景 |
|------|--------|----------|----------|
| Vue CDN版 | ⭐ | ⭐ | 简单页面、快速原型 |
| Vue + Vite | ⭐⭐⭐ | ⭐⭐⭐ | 正式项目、需要构建 |
| Vue + Tailwind | ⭐⭐⭐⭐ | ⭐⭐⭐ | 需要复杂样式 |

**推荐方案：Vue CDN版（零配置）**

## 2. 项目结构（极简）

**单文件方案（推荐）：**
```
/
├── index.html          # 包含所有代码的单一文件
├── product.json        # 产品数据（可选）
└── assets/            # 图片视频资源
    ├── images/
    └── videos/
```

**多文件方案（适度扩展）：**
```
/
├── index.html          # 主入口
├── app.js             # Vue应用逻辑
├── style.css          # 样式表
├── data.js            # 产品数据
└── assets/            # 媒体资源
```

## 3. 部署方式

### 3.1 极简部署（无需构建）

**index.html完整代码：**
```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>农产品溯源系统</title>
  <script src="https://unpkg.com/vue@3/dist/vue.global.js"></script>
  <style>
    /* 所有样式 */
  </style>
</head>
<body>
  <div id="app">
    <!-- 页面内容 -->
  </div>
  
  <script>
    const { createApp, ref, computed } = Vue;
    
    createApp({
      setup() {
        // 业务逻辑
        return {};
      }
    }).mount('#app');
  </script>
</body>
</html>
```

**部署步骤：**
1. 将index.html上传到任意静态服务器
2. 访问URL即可运行

**支持的部署平台：**
- GitHub Pages（免费）
- Vercel/Netlify（免费）
- 阿里云OSS（按量付费）
- 腾讯云COS（按量付费）
- Nginx/Apache服务器
- 甚至可直接双击打开HTML文件

### 3.2 单命令部署

**使用Vite（需要Node.js）：**
```bash
# 初始化项目
npm create vite@latest traceability -- --template vue
cd traceability
npm install
npm run dev

# 构建（生成dist目录）
npm run build

# 部署dist目录到任意静态服务器
```

## 4. 数据管理

**静态数据方案（最简）：**
```javascript
const productData = {
  id: '4395',
  name: '枣甜5号',
  code: '4395',
  location: '新疆维吾尔自治区和田地区',
  plantingDate: '2024-03-15',
  images: [
    'https://example.com/image1.jpg',
    'https://example.com/image2.jpg'
  ],
  videos: [
    'https://example.com/video.mp4'
  ],
  quality: {
    harvestStart: '2024-09-01',
    harvestEnd: '2024-09-15',
    sweetness: 22.5,
    weight: 15.3,
    texture: '肉质脆嫩，甜度高',
    crowd: ['一般人群', '糖尿病患者慎用'],
    summary: '今年光照充足，糖度较往年提升15%，品质优良。'
  }
};
```

## 5. 组件设计（内联方案）

**不需要独立组件文件，直接在index.html中实现：**

```html
<div id="app">
  <!-- Hero区域 -->
  <section class="hero">
    <img :src="product.videos[0]" class="hero-video">
    <div class="hero-overlay">
      <span class="badge">{{ product.code }}</span>
      <h1>{{ product.name }}</h1>
    </div>
  </section>
  
  <!-- 产品信息 -->
  <section class="info-cards">
    <div class="card">
      <h3>品种名称</h3>
      <p>{{ product.name }}</p>
    </div>
    <!-- 其他信息卡片 -->
  </section>
  
  <!-- 媒体展示 -->
  <section class="media-gallery">
    <div class="image-grid">
      <img v-for="img in product.images" :src="img">
    </div>
    <video :src="product.videos[0]"></video>
  </section>
  
  <!-- 质量信息 -->
  <section class="quality-info">
    <div class="quality-card">
      <span class="label">糖度</span>
      <span class="value">{{ product.quality.sweetness }}°Brix</span>
    </div>
    <!-- 其他质量指标 -->
  </section>
</div>
```

## 6. 样式规范（纯CSS）

**关键样式变量：**
```css
:root {
  --primary: #4A7C59;
  --secondary: #F5E6D3;
  --accent: #E8A838;
  --neutral: #2D3436;
  --bg: #FAFAFA;
  --card-radius: 16px;
  --card-shadow: 0 4px 16px rgba(0,0,0,0.12);
}
```

**响应式断点（极简）：**
```css
/* 移动端（默认） */
.container { width: 100%; padding: 0 16px; }

/* 平板及以上 */
@media (min-width: 768px) {
  .container { width: 720px; margin: 0 auto; }
}

/* 桌面端 */
@media (min-width: 1024px) {
  .container { width: 960px; margin: 0 auto; }
}
```

## 7. 部署检查清单

**上线前检查：**
- [ ] 图片和视频路径正确（使用绝对URL或相对路径）
- [ ] 所有外部资源（CDN链接）可访问
- [ ] 移动端布局正常
- [ ] 视频播放控件正常
- [ ] 页面加载速度可接受（<3秒）

**部署平台选择：**
1. **GitHub Pages**（推荐，免费）
   - 上传index.html到仓库
   - Settings → Pages → 启用
   
2. **本地Nginx/Apache**
   - 直接放置index.html
   - 配置SSL证书（可选）

3. **云对象存储**
   - 阿里云OSS
   - 腾讯云COS
   - 上传文件，配置域名访问

## 8. 性能优化（可选）

**图片优化：**
```html
<!-- 懒加载 -->
<img loading="lazy" src="image.jpg">

<!-- 响应式图片 -->
<img srcset="small.jpg 480w, large.jpg 1080w">
```

**视频优化：**
```html
<!-- 添加海报 -->
<video poster="thumbnail.jpg" src="video.mp4"></video>

<!-- 延迟加载 -->
<video loading="lazy" src="video.mp4"></video>
```

## 9. 快速开始模板

**完整index.html模板：**
```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>农产品溯源 - 枣甜5号</title>
  <style>
    /* 基础重置 */
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body { font-family: 'Noto Sans SC', sans-serif; background: var(--bg); }
    
    /* 变量 */
    :root {
      --primary: #4A7C59;
      --secondary: #F5E6D3;
      --accent: #E8A838;
      --neutral: #2D3436;
      --bg: #FAFAFA;
      --card-radius: 16px;
      --card-shadow: 0 4px 16px rgba(0,0,0,0.12);
    }
    
    /* 布局 */
    .container { max-width: 960px; margin: 0 auto; padding: 20px; }
    
    /* 卡片 */
    .card {
      background: white;
      border-radius: var(--card-radius);
      padding: 24px;
      box-shadow: var(--card-shadow);
      margin-bottom: 20px;
    }
    
    /* 图片 */
    .image-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
      gap: 16px;
    }
    .image-grid img {
      width: 100%;
      border-radius: var(--card-radius);
      object-fit: cover;
    }
    
    /* 视频 */
    video {
      width: 100%;
      border-radius: var(--card-radius);
      margin-top: 20px;
    }
    
    /* 标签 */
    .badge {
      display: inline-block;
      background: var(--accent);
      color: white;
      padding: 4px 12px;
      border-radius: 9999px;
      font-size: 14px;
    }
  </style>
</head>
<body>
  <div id="app">
    <div class="container">
      <!-- 产品信息 -->
      <div class="card">
        <span class="badge">{{ product.code }}</span>
        <h1>{{ product.name }}</h1>
        <p>📍 {{ product.location }}</p>
        <p>🌱 定植时间：{{ product.plantingDate }}</p>
      </div>
      
      <!-- 图片展示 -->
      <div class="card">
        <h2>产品图片</h2>
        <div class="image-grid">
          <img v-for="img in product.images" :src="img" loading="lazy">
        </div>
      </div>
      
      <!-- 视频展示 -->
      <div class="card">
        <h2>产品视频</h2>
        <video controls :src="product.videos[0]" poster="thumbnail.jpg"></video>
      </div>
      
      <!-- 质量信息 -->
      <div class="card">
        <h2>采收质量</h2>
        <p>糖度：{{ product.quality.sweetness }}°Brix</p>
        <p>重量：{{ product.quality.weight }}g</p>
        <p>口感：{{ product.quality.texture }}</p>
        <p>人群：{{ product.quality.crowd.join('、') }}</p>
        <blockquote>{{ product.quality.summary }}</blockquote>
      </div>
    </div>
  </div>

  <script src="https://unpkg.com/vue@3/dist/vue.global.js"></script>
  <script>
    const { createApp, ref } = Vue;
    
    createApp({
      setup() {
        const product = ref({
          id: '4395',
          name: '枣甜5号',
          code: '4395',
          location: '新疆维吾尔自治区和田地区',
          plantingDate: '2024-03-15',
          images: [
            'https://picsum.photos/400/300?random=1',
            'https://picsum.photos/400/300?random=2',
            'https://picsum.photos/400/300?random=3'
          ],
          videos: [
            'https://www.w3schools.com/html/mov_bbb.mp4'
          ],
          quality: {
            harvestStart: '2024-09-01',
            harvestEnd: '2024-09-15',
            sweetness: 22.5,
            weight: 15.3,
            texture: '肉质脆嫩，甜度高，入口即化',
            crowd: ['一般人群', '注重健康饮食者'],
            summary: '得益于和田地区充足的日照和昼夜温差，今年的枣甜5号品质优于往年，糖度达到22.5°Brix，属于特级品质。'
          }
        });
        
        return { product };
      }
    }).mount('#app');
  </script>
</body>
</html>
```

**部署这个文件：**
1. 保存为 `index.html`
2. 上传到GitHub Pages/Vercel/任意服务器
3. 访问即可运行
