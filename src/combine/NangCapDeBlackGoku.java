/*
 * Click nbfs://nbhost/SystemFileSystem/Templates/Licenses/license-default.txt to change this license
 * Click nbfs://nbhost/SystemFileSystem/Templates/Classes/Class.java to edit this template
 */
package combine;

import consts.ConstNpc;
import item.Item;
import player.Player;
import services.ItemService;
import services.Service;
import player.Service.InventoryService;
import player.system.Template;
import utils.Util;

/**
 *
 * @author Administrator
 */
public class NangCapDeBlackGoku {
    private static final int GOLD_TAO_DA = 50_000_000;
    private static final int RATIO_TAO_DA = 80;

    public static void showInfoCombine(Player player) {
        if (player.combineNew.itemsCombine.size() == 1) {
            Item DuiDuc = player.combineNew.itemsCombine.get(0);
            if (DuiDuc.template.id == 568 && DuiDuc.quantity >= 15) {
                player.combineNew.goldCombine = GOLD_TAO_DA;
                player.combineNew.ratioCombine = RATIO_TAO_DA;

                String npcSay = "|2|Tạo Đệ Black Từ Trứng Mabư\n";
                npcSay += "|2|Cần 15 trứng mabư\n";
                npcSay += "|2|Tỉ lệ thành công: " + player.combineNew.ratioCombine + "%\n";
                npcSay += "|2|Cần: " + Util.numberToMoney(player.combineNew.goldCombine) + " vàng\n";
                npcSay += "|7|Thất bại -15 trứng mabư\n";
                if (player.inventory.gold < player.combineNew.goldCombine) {
                    npcSay += "|7|Còn thiếu " + Util.powerToString(player.combineNew.goldCombine - player.inventory.gold) + " vàng\n";
                    CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU, npcSay, "Đóng");
                } else {
                    CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.MENU_START_COMBINE, npcSay,
                            "Tạo Đệ Black Goku\n" + Util.numberToMoney(player.combineNew.goldCombine) + " vàng\n", "Từ chối");
                }
            } else {
                CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                        "Cần 15 trứng mabư", "Đóng");
            }
        } else {
            CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                    "Cần 15 trứng mabư", "Đóng");
        }
    }

    public static void NangCapDeBlackGoku(Player player) {
        if (player.combineNew.itemsCombine.size() == 1) {
            int gold = player.combineNew.goldCombine;

            if (player.inventory.gold < gold) {
                Service.gI().sendThongBao(player, "Không đủ vàng để thực hiện");
                return;
            }

            Item TrungMabu = player.combineNew.itemsCombine.get(0);
            if (TrungMabu.template.id == 568 && TrungMabu.quantity >= 15) {
                player.inventory.gold -= gold;
                InventoryService.gI().subQuantityItemsBag(player, TrungMabu, 15);
                if (Util.isTrue(player.combineNew.ratioCombine, 100)) {
                    Template.ItemTemplate DeTuBlackGoku = ItemService.gI().getTemplate(1774); 
                    Item DeBlackGoku = new Item();
                    DeBlackGoku.template = DeTuBlackGoku;
                    DeBlackGoku.quantity = 1;
                    InventoryService.gI().addItemBag(player, DeBlackGoku);
                    CombineService.gI().sendEffectSuccessCombine(player);
                } else {
                    CombineService.gI().sendEffectFailCombine(player);
                }

                InventoryService.gI().sendItemBags(player);
                Service.gI().sendMoney(player);
                CombineService.gI().reOpenItemCombine(player);
            } else {
                Service.gI().sendThongBao(player, "Không đủ trứng mabư");
            }
        }
    }
}
