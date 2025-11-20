package combine;

import consts.ConstNpc;
import item.Item;
import player.Player;
import player.Service.InventoryService;
import services.ItemService;
import services.Service;
import player.system.Template;
import utils.Util;

public class NangCapItemCap2 {

    private static final int GOLD_TAO_DA = 50_000_000;  // Vàng cần thiết để tạo đá Hematite
    private static final int RATIO_TAO_DA = 80;  // Tỉ lệ thành công tạo đá Hematite (80%)

    public static void showInfoCombine(Player player) {
        if (player.combineNew.itemsCombine.size() == 1) {
            Item itemc1 = player.combineNew.itemsCombine.get(0);
            if (itemc1.template.id >= 381 && itemc1.template.id <= 385 && itemc1.quantity >= 10) {
                player.combineNew.goldCombine = GOLD_TAO_DA;
                player.combineNew.ratioCombine = RATIO_TAO_DA;

                String npcSay = "|2|Tạo đá Item c2 từ Item Cấp 1\n";
                npcSay += "|2|Cần 10 item cấp 1 lên item cấp 2\n";
                npcSay += "|2|Tỉ lệ thành công: " + player.combineNew.ratioCombine + "%\n";
                npcSay += "|2|Cần: " + Util.numberToMoney(player.combineNew.goldCombine) + " vàng\n";
                npcSay += "|7|Thất bại -10 item cấp 1\n";

                // Kiểm tra tài nguyên và đưa ra menu
                if (player.inventory.gold < player.combineNew.goldCombine) {
                    npcSay += "|7|Còn thiếu " + Util.powerToString(player.combineNew.goldCombine - player.inventory.gold) + " vàng\n";
                    CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU, npcSay, "Đóng");
                } else {
                    CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.MENU_START_COMBINE, npcSay,
                            "Tạo item c2\n" + Util.numberToMoney(player.combineNew.goldCombine) + " vàng\n", "Từ chối");
                }
            } else {
                CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                        "Cần 10 item Cấp 1", "Đóng");
            }
        } else {
            CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                    "Cần 10 item cấp 1", "Đóng");
        }
    }

    public static void Itemc2(Player player) {
        if (player.combineNew.itemsCombine.size() == 1) {
            int gold = player.combineNew.goldCombine;

            if (player.inventory.gold < gold) {
                Service.gI().sendThongBao(player, "Không đủ vàng để thực hiện");
                return;
            }

            Item itemc1 = player.combineNew.itemsCombine.get(0);

            // Kiểm tra xem có đủ sao pha lê không
            if (itemc1.template.id >= 381 && itemc1.template.id <= 385 && itemc1.quantity >= 10) {
                player.inventory.gold -= gold;
                InventoryService.gI().subQuantityItemsBag(player, itemc1, 10);
                if (Util.isTrue(player.combineNew.ratioCombine, 100)) {
                    int randomId = Util.nextInt(1150,1154); // Sử dụng hàm random trong game server của bạn
                    Template.ItemTemplate Itemc2Template = ItemService.gI().getTemplate(randomId);
                    Item itemc2 = new Item();
                    itemc2.template = Itemc2Template;
                    itemc2.quantity = 1;
                    InventoryService.gI().addItemBag(player, itemc2);
                    CombineService.gI().sendEffectSuccessCombine(player);
                } else {
                    CombineService.gI().sendEffectFailCombine(player);
                }

                InventoryService.gI().sendItemBags(player);
                Service.gI().sendMoney(player);
                CombineService.gI().reOpenItemCombine(player);
            } else {
                Service.gI().sendThongBao(player, "Không đủ item c1 tạo item c2");
            }
        }
    }
}
