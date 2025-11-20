/*
 * Click nbfs://nbhost/SystemFileSystem/Templates/Licenses/license-default.txt to change this license
 * Click nbfs://nbhost/SystemFileSystem/Templates/Classes/Class.java to edit this template
 */
package combine;

import consts.ConstNpc;
import item.Item;
import map.Service.ChangeMapService;
import player.Player;
import services.ItemService;
import services.Service;
import player.Service.InventoryService;
import player.system.Template;
import services.PetService;
import utils.Util;

/**
 *
 * @author Administrator
 */
public class NangCapDeTuBlackGokuRose {
    private static final int GOLD_TAO_DA = 500_000_000;

    public static void showInfoCombine(Player player) {
        if (player.combineNew.itemsCombine.size() == 1) {
            Item Damai = player.combineNew.itemsCombine.get(0);
            if (Damai.template.id == 1439 && Damai.quantity >= 1) {
                player.combineNew.goldCombine = GOLD_TAO_DA;

                String npcSay = "|2|nâng cấp đệ black\n";
                npcSay += "tháo đồ đệ tử trước khi nâng cấp đệ Black Goku Rose\nkhông sẽ mất bên admin không chịu trách nghiệm\n";
                npcSay += "|2|Tỉ lệ thành công: 1% - 100%\n";
                npcSay += "|2|Cần: " + Util.numberToMoney(player.combineNew.goldCombine) + " vàng\n";
                npcSay += "|7|Thất bại - 500tr vàng và 1 đá mài\n";
                if (player.inventory.gold < player.combineNew.goldCombine) {
                    npcSay += "|7|Còn thiếu " + Util.powerToString(player.combineNew.goldCombine - player.inventory.gold) + " vàng\n và";
                    CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU, npcSay, "Đóng");
                } else {
                    CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.MENU_START_COMBINE, npcSay,
                            "Nâng Đệ Black Goku Rose\n" + Util.numberToMoney(player.combineNew.goldCombine) + " vàng\n", "Từ chối");
                }
            } else {
                CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                        "Cần 500tr vàng", "Đóng");
            }
        } else {
            CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                    "Cần 1 đệ tử black goku", "Đóng");
            }
        }

    public static void NangCapDeBlackGokuRose(Player player) {
        if (player.combineNew.itemsCombine.size() == 1) {
            int gold = player.combineNew.goldCombine;

            if (player.inventory.gold < gold) {
                Service.gI().sendThongBao(player, "Không đủ vàng để thực hiện");
                return;
            }
            if (player.pet == null || player.pet.typePet != 3) {
                Service.gI().sendThongBao(player, "Ngươi không có đệ tử Black goku");
                return;
            }
            Item Damai = player.combineNew.itemsCombine.get(0);
            if (Damai.template.id == 1439 && Damai.quantity >= 1) {
                player.inventory.gold -= gold;
                if (Util.isTrue(10, 50)) {
                    ChangeMapService.gI().exitMap(player.pet);
                    PetService.gI().createBlackGokuRose(player, player.gender);
                    CombineService.gI().sendEffectSuccessCombine(player);
                } else {
                    CombineService.gI().sendEffectFailCombine(player);
                }
                InventoryService.gI().subQuantityItemsBag(player, Damai,1);
                InventoryService.gI().sendItemBags(player);
                Service.gI().sendMoney(player);
                CombineService.gI().reOpenItemCombine(player);
            } else {
                Service.gI().sendThongBao(player, "Không đủ đá mài");
            }
        }
    }
}
