<div align="center">
  <div align="center">
    
  </div>

# PTT to Discord

  <div>

  </div>

  <p align="center"><font>使您追蹤的PTT貼文即時轉發至 Discord</font></p>
</div>

---
> [!IMPORTANT]\
> 本專案僅提供瀏覽文章留言功能，無意違反任何 PTT 站方規範。如有牴觸，請馬上告知。

> PTT 無法於雲端伺服器登入，故請於本地或其他環境中執行。

---

專案結構:

```
app            主程式
app/ptd        主程式-進入點
app/ent        主程式-資料庫
app/config     主程式-設定

core           功能核心
core/event     內建事件

uao            Big5-UAO <-> UTF-8 轉換
```