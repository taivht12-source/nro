package npc.list;

/*
 * @Author Coder: Nguyễn Tấn Tài
 * @Description: Ngọc Rồng Kiwi - Máy Chủ Chuẩn Teamobi 2025
 * @Group Zalo: https://zalo.me/g/toiyeuvietnam2025
 */
import consts.ConstNpc;
import database.PlayerDAO;
import item.Item;
import npc.Npc;
import player.Player;
import player.Service.InventoryService;
import services.ItemService;
import services.PetService;
import services.Service;
import services.TaskService;
import services.func.Input;
import shop.ShopService;
import utils.Util;

public class OngGohan extends Npc {

    public OngGohan(int mapId, int status, int cx, int cy, int tempId, int avartar) {
        super(mapId, status, cx, cy, tempId, avartar);
    }

    @Override
    public void openBaseMenu(Player player) {
        if (canOpenNpc(player)) {
            if (!TaskService.gI().checkDoneTaskTalkNpc(player, this)) {
                this.createOtherMenu(player, ConstNpc.BASE_MENU,
                        "|0| BlackGoku - Game Ngọc Rồng Chuẩn Teamobi 2025",
                        "Giftcode",
                        "Đổi Mật Khẩu",
                        "Nhận ngọc xanh",
                        "Shop Đổi skill\nĐệ tử");
            }
        }
    }

    @Override
    public void confirmMenu(Player player, int select) {
        if (canOpenNpc(player)) {
            if (player.idMark.isBaseMenu()) {
                switch (select) {
                    case 0 -> Input.gI().createFormGiftCode(player);
                    case 1 -> Input.gI().createFormChangePassword(player);
                    case 2 -> {
                    if (player.inventory.gem >= 20_000_000) {
                    this.npcChat(player, "Bạn đã ngọc xanh nhận rồi");
                    break;
                    }
                    player.inventory.gem += 2000000;
                    Service.gI().sendMoney(player);
                    Service.gI().sendThongBao(player, "Bạn vừa nhận được 2 triệu ngọc xanh!");
                    }
                    case 3 -> {
                    ShopService.gI().opendShop(player, "DOI_SKILL_DE", false);
                    }
                }
            }
        }
    }
}