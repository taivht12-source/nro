package combine;

import consts.ConstNpc;
import item.Item;
import player.Player;
import player.Service.InventoryService;
import services.ItemService;
import services.Service;
import player.system.Template;
import utils.Util;

public class DanhBongSaoPhaLe {

    private static final int GOLD_NANG_CAP = 100_000_000;
    private static final int RATIO_NANG_CAP = 100;

    // Hiển thị thông tin nâng cấp Sao Pha Lê
    public static void showInfoCombine(Player player) {
        if (player.combineNew.itemsCombine.size() == 2) {
            Item saoPhaLe = null;
            Item daMai = null;

            // Tìm kiếm sao Pha Lê và đá mài
            for (Item item : player.combineNew.itemsCombine) {
                if (item.template.id >= 1416 && item.template.id <= 1422) {
                    saoPhaLe = item;
                } else if (item.template.id == 1439) {
                    daMai = item;
                }
            }

            if (saoPhaLe != null && daMai != null && saoPhaLe.quantity >= 2) {
                player.combineNew.goldCombine = GOLD_NANG_CAP;
                player.combineNew.ratioCombine = RATIO_NANG_CAP;

                String npcSay = "|2|Nâng cấp Sao Pha Lê từ cấp 2 lên Sao Pha Lê lấp lánh\n";
                npcSay += "|2|Tỉ lệ thành công: " + player.combineNew.ratioCombine + "%\n";
                npcSay += "|2|Cần 1 đá mài\n";
                npcSay += "|2|Cần: " + Util.numberToMoney(player.combineNew.goldCombine) + " vàng\n";
                npcSay += "|7|Thất bại -1 đá đá mài\n";

                // Kiểm tra tài nguyên và đưa ra menu
                if (player.inventory.gold < player.combineNew.goldCombine) {
                    npcSay += "|7|Còn thiếu " + Util.powerToString(player.combineNew.goldCombine - player.inventory.gold) + " vàng\n";
                    CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU, npcSay, "Đóng");
                } else {
                    CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.MENU_START_COMBINE, npcSay,
                            "Nâng cấp\n" + Util.numberToMoney(player.combineNew.goldCombine) + " vàng\n"
                            + Util.numberToMoney(player.combineNew.gemCombine) + " ngọc\n", "Từ chối");
                }
            } else {
                CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                        "Cần x2 Sao Pha Lê cấp 2 và 1 đá mài", "Đóng");
            }
        } else {
            CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                    "Cần x2 Sao Pha Lê cấp 2 và 1 đá mài", "Đóng");
        }
    }

    // Thực hiện nâng cấp Sao Pha Lê
    public static void danhBongSaoPhaLe(Player player) {
        if (player.combineNew.itemsCombine.size() == 2) {
            int gold = player.combineNew.goldCombine;
            int gem = player.combineNew.gemCombine;

            // Kiểm tra đủ vàng và ngọc để thực hiện nâng cấp
            if (player.inventory.gold < gold) {
                Service.gI().sendThongBao(player, "Không đủ vàng để thực hiện");
                return;
            }
            if (player.inventory.gem < gem) {
                Service.gI().sendThongBao(player, "Không đủ ngọc để thực hiện");
                return;
            }

            Item saoPhaLe = null;
            Item daMai = null;

            // Tìm kiếm sao Pha Lê và đá mài
            for (Item item : player.combineNew.itemsCombine) {
                if (item.template.id >= 1416 && item.template.id <= 1422) {
                    saoPhaLe = item;
                } else if (item.template.id == 1439) {
                    daMai = item;
                }
            }

            if (saoPhaLe != null && daMai != null && saoPhaLe.quantity >= 2) {
                player.inventory.gold -= gold;
                player.inventory.gem -= gem;

                if (Util.isTrue(player.combineNew.ratioCombine, 100)) {
                    // Tính toán ID của Sao Pha Lê lấp lánh
                    int saoPhaLeLapLanhId = 1426 + (saoPhaLe.template.id - 1416);
                    Template.ItemTemplate newTemplate = ItemService.gI().getTemplate(saoPhaLeLapLanhId);
                    Item newItem = new Item();
                    newItem.template = newTemplate;
                    newItem.quantity = 1;

                    // Sao chép các Option từ Sao Pha Lê cấp 2 vào Sao Pha Lê cấp 3 và thêm +1 vào option
                    for (Item.ItemOption option : saoPhaLe.itemOptions) {
                        // Thêm 1 vào giá trị option
                        Item.ItemOption newOption = new Item.ItemOption(option.optionTemplate.id, option.param + 1);
                        newItem.itemOptions.add(newOption);
                    }

                    // Thêm Sao Pha Lê cấp 3 vào túi đồ của người chơi
                    InventoryService.gI().addItemBag(player, newItem);
                    InventoryService.gI().subQuantityItemsBag(player, saoPhaLe, 2); // Trừ 2 Sao Pha Lê cấp 2
                    CombineService.gI().sendEffectSuccessCombine(player); // Hiển thị hiệu ứng thành công
                } else {
                    CombineService.gI().sendEffectFailCombine(player); // Hiển thị hiệu ứng thất bại
                }

                InventoryService.gI().subQuantityItemsBag(player, daMai, 1); // Tiêu thụ đá mài
                InventoryService.gI().sendItemBags(player); // Cập nhật túi đồ
                Service.gI().sendMoney(player); // Cập nhật tiền
                CombineService.gI().reOpenItemCombine(player); // Mở lại giao diện nâng cấp
            }
        }
    }

    // Hàm chuyển đổi ID của Sao Pha Lê cấp 1 thành ID Sao Pha Lê cấp 3
    private static int saoPhaLeLapLanhId(int saoPhaLeCap1Id) {
        switch (saoPhaLeCap1Id) {
            case 1416:
                return 1426; // Sao pha lê đỏ
            case 1417:
                return 1427; // Sao pha lê lam
            case 1418:
                return 1428; // Sao pha lê hồng
            case 1419:
                return 1429; // Sao pha lê tím
            case 1420:
                return 1430; // Sao pha lê cam
            case 1421:
                return 1431; // Sao pha lê vàng
            case 1422:
                return 1432; // Sao pha lê lục
            default:
                return -1; // Trường hợp không hợp lệ
        }
    }
}
