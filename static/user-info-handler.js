// User Info Handler for ESG Dashboard
(function() {
    "use strict";

    // Auto-refresh userInfo from backend to keep DID and other fields up-to-date
    function refreshUserInfoFromBackend() {
        const userInfoStr = localStorage.getItem("userInfo");
        if (!userInfoStr) return;
        try {
            const userInfo = JSON.parse(userInfoStr);
            const email = userInfo.Email || localStorage.getItem("esg_email") || "";
            if (!email) return;
            fetch("/api/register/user-detail?email=" + encodeURIComponent(email))
            .then(function(r){return r.json()})
            .then(function(j){
                if(j.code===200 && j.data){
                    const d = j.data;
                    const updated = {
                        Email: d.email || email,
                        FullName: d.full_name || d.fullname || email,
                        Name: d.full_name || d.fullname || email,
                        Role: d.role || "user",
                        Phone: d.phone || "",
                        Age: d.age || 20,
                        DID: d.did || "",
                        RegisterTime: d.created_at || "",
                        BuildingName: d.building_name || "",
                        BuildingAddr: d.building_addr || "",
                        BuildingType: d.building_type || "",
                        PropertyName: d.property_name || "",
                        Occupation: d.occupation || "",
                        Institution: d.institution || "",
                        isLoggedIn: true
                    };
                    localStorage.setItem("userInfo", JSON.stringify(updated));
                    console.log("[UserInfo] Refreshed from backend, DID:", updated.DID);
                }
            }).catch(function(){});
        } catch(e) {}
    }
    refreshUserInfoFromBackend();


    function initUserInfoDisplay() {
        console.log("[UserInfo] Starting...");

        const userInfoStr = localStorage.getItem("userInfo");
        if (!userInfoStr) {
            console.log("[UserInfo] No userInfo found");
            return;
        }

        try {
            const userInfo = JSON.parse(userInfoStr);
            console.log("[UserInfo] Data:", userInfo);

            const emButton = document.querySelector(".w-10.h-10.rounded-full");
            if (!emButton) {
                console.log("[UserInfo] Button not found");
                return;
            }

            let initials = "EM";
            const esgFullName = localStorage.getItem("esg_fullName");
            let displayName = esgFullName && !esgFullName.includes("@") ? esgFullName : (userInfo.Email || "U");
            const nameParts = displayName.split(" ");
            initials = nameParts.length >= 2 ? (nameParts[0][0] + nameParts[1][0]).toUpperCase() : displayName.substring(0, 2).toUpperCase();

            emButton.textContent = initials;
            emButton.onclick = function() {
                    var email = localStorage.getItem('esg_email') || (JSON.parse(localStorage.getItem('userInfo')||'{}').Email) || '';
                    if(!email){ var u=JSON.parse(localStorage.getItem('userInfo')||'{}'); if(u.Email)email=u.Email; }
                    if(email){
                        fetch('/api/register/user-detail?email='+encodeURIComponent(email)).then(function(r){return r.json()}).then(function(j){
                            if(j.code===200 && j.data){
                                var d=j.data;
                                var info={Email:d.email||email,FullName:d.full_name||'',Name:d.full_name||'',Role:d.role||'',Phone:d.phone||'',Age:d.age||'',DID:d.did||'',RegisterTime:d.created_at||'',BuildingName:d.building_name||'',BuildingAddr:d.building_addr||'',BuildingType:d.building_type||'',PropertyName:d.property_name||'',Occupation:d.occupation||'',Institution:d.institution||'',isLoggedIn:true};
                                localStorage.setItem('userInfo',JSON.stringify(info));
                                showUserInfoModal(info);
                            } else { var f2=localStorage.getItem('userInfo'); if(f2)showUserInfoModal(JSON.parse(f2)); }
                        }).catch(function(){ var f2=localStorage.getItem('userInfo'); if(f2)showUserInfoModal(JSON.parse(f2)); });
                    } else { var f2=localStorage.getItem('userInfo'); if(f2)showUserInfoModal(JSON.parse(f2)); }
                };
            console.log("[UserInfo] Initialized with initials:", initials);
        } catch(e) {
            console.error("[UserInfo] Error:", e);
        }
    }

    function showUserInfoModal(userInfo) {
        const modal = document.createElement("div");
        modal.style.cssText = "position:fixed;top:0;left:0;width:100%;height:100%;background:rgba(0,0,0,0.5);z-index:9999;padding:20px;display:flex;align-items:center;overflow-y:auto;";

        const modalContent = document.createElement("div");
        modalContent.style.cssText = "background:white;border-radius:12px;max-width:600px;width:100%;margin:auto;box-shadow:0 20px 25px -5px rgba(0,0,0,0.1);max-height:95vh;overflow-y:auto;";

        const header = document.createElement("div");
        header.style.cssText = "padding:20px;text-align:center;border-bottom:1px solid #e2e8f0;";
        header.innerHTML = "<h2 style='margin:0;display:flex;align-items:center;justify-content:center;gap:10px;color:#1e293b;'><span style='width:40px;height:40px;background:linear-gradient(135deg,#667eea,#764ba2);border-radius:50%;display:flex;align-items:center;justify-content:center;color:white;font-weight:bold;'>" + getInitials(userInfo) + "</span>User Information</span></h2>";

        const content = document.createElement("div");
        content.style.cssText = "padding:20px;";

        const fields = [
            {key:"FullName",label:"Full Name"},{key:"Name",label:"Name"},{key:"Email",label:"Email Address"},
            {key:"Phone",label:"Phone Number"},{key:"DID",label:"DID"},{key:"Role",label:"Role"},{key:"Age",label:"Age"},
            {key:"RegisterTime",label:"Registration Time"},{key:"BuildingName",label:"Building Name"},
            {key:"BuildingAddr",label:"Building Address"},{key:"BuildingType",label:"Building Type"},
            {key:"PropertyName",label:"Property Name"},{key:"Occupation",label:"Occupation"},{key:"Institution",label:"Institution"}
        ];

        let html = "";
        let count = 0;
        fields.forEach(f => {
            let val = userInfo[f.key];
            if(f.key === "FullName" && localStorage.getItem("esg_fullName") && !localStorage.getItem("esg_fullName").includes("@")) {
                val = localStorage.getItem("esg_fullName");
            }
            if(val) {
                count++;
                html += "<div style='margin-bottom:16px;'><label style='display:block;font-size:11px;color:#64748b;font-weight:600;text-transform:uppercase;margin-bottom:6px;'>" + f.label + "</label><div style='font-size:14px;color:#1e293b;padding:12px;background-color:#f8fafc;border-radius:6px;border:1px solid #e2e8f0;word-wrap:break-word;'>" + escapeHtml(val) + "</div></div>";
            }
        });
        console.log("[UserInfo] Showing " + count + " fields");

        const footer = document.createElement("div");
        footer.style.cssText = "padding:20px;border-top:1px solid #e2e8f0;display:flex;gap:10px;";
        footer.innerHTML = "<button id='close-btn' style='flex:1;padding:12px 20px;background:#e2e8f0;color:#1e293b;border:none;border-radius:8px;font-weight:600;cursor:pointer;'>Close</button><button id='logout-btn' style='flex:1;padding:12px 20px;background:#ef4444;color:white;border:none;border-radius:8px;font-weight:600;cursor:pointer;'>Logout</button>";

        modalContent.appendChild(header);
        content.innerHTML = html;
        modalContent.appendChild(content);
        modalContent.appendChild(footer);
        modal.appendChild(modalContent);
        document.body.appendChild(modal);

        document.getElementById("close-btn").onclick = () => modal.remove();
        document.getElementById("logout-btn").onclick = () => { localStorage.removeItem("userInfo");localStorage.removeItem("userToken");localStorage.removeItem("isLoggedIn");window.location.href="/static/index.html#login"; };
        modal.onclick = (e) => { if(e.target === modal) modal.remove(); };
        console.log("[UserInfo] Modal shown");
    }

    function getInitials(userInfo) {
        const esgFullName = localStorage.getItem("esg_fullName");
        let name = esgFullName && !esgFullName.includes("@") ? esgFullName : (userInfo.Email || "U");
        const parts = name.split(" ");
        return parts.length >= 2 ? (parts[0][0] + parts[1][0]).toUpperCase() : name.substring(0, 2).toUpperCase();
    }

    function escapeHtml(t) {
        const d = document.createElement("div");
        d.textContent = t;
        return d.innerHTML;
    }

    function waitForReactApp() {
        let tries = 0;
        const interval = setInterval(() => {
            tries++;
            const btn = document.querySelector(".w-10.h-10.rounded-full");
            if(btn) {
                console.log("[UserInfo] Button found!");
                clearInterval(interval);
                initUserInfoDisplay();
            } else if(tries >= 10) {
                clearInterval(interval);
            }
        }, 500);
    }

    if(document.readyState === "loading") {
        document.addEventListener("DOMContentLoaded", waitForReactApp);
    } else {
        waitForReactApp();
    }

    console.log("[UserInfo] Script loaded");
})();
