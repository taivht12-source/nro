package combine;

/*
 * @Author Coder: Nguyễn Tấn Tài
 * @Description: Ngọc Rồng Kiwi - Máy Chủ Chuẩn Teamobi 2025
 * @Group Zalo: https://zalo.me/g/toiyeuvietnam2025
 */
import consts.ConstNpc;
import item.Item;
import player.Player;
import player.Service.InventoryService;
import services.ItemService;
import services.Service;
import player.system.Template.ItemTemplate;
import utils.Util;

public class NangCapSaoPhaLe {

    private static final int GOLD_NANG_CAP = 200_000_000;
    private static final int GEM_NANG_CAP = 10;
    private static final int RATIO_NANG_CAP = 50;

    // Hiển thị thông tin nâng cấp sao Pha Lê
    public static void showInfoCombine(Player player) {
        if (player.combineNew.itemsCombine.size() == 2) {
            Item saoPhaLe = null;
            Item hematite = null;

            // Tìm kiếm sao Pha Lê cấp 1 và Hematite
            for (Item item : player.combineNew.itemsCombine) {
                if (item.template.id >= 441 && item.template.id <= 447) {
                    saoPhaLe = item; // Lấy sao pha lê cấp 1
                } else if (item.template.id == 1423) {
                    hematite = item; // Lấy đá hematite
                }
            }

            if (saoPhaLe != null && hematite != null) {
                player.combineNew.goldCombine = GOLD_NANG_CAP;
                player.combineNew.gemCombine = GEM_NANG_CAP;
                player.combineNew.ratioCombine = RATIO_NANG_CAP;

                String npcSay = "|2|Nâng cấp Sao Pha Lê từ cấp 1 lên cấp 2\n";
                npcSay += "|2|Tỉ lệ thành công: " + player.combineNew.ratioCombine + "%\n";
                npcSay += "|2|Cần 1 đá Hematite\n";
                npcSay += "|2|Cần: " + Util.numberToMoney(player.combineNew.gemCombine) + " ngọc\n";
                npcSay += "|2|Cần: " + Util.numberToMoney(player.combineNew.goldCombine) + " vàng\n";
                npcSay += "|7|Thất bại -1 đá Hematite\n";

                // Kiểm tra tài nguyên và đưa ra menu
                if (player.inventory.getGem() < player.combineNew.gemCombine) {
                    npcSay += "|7|Còn thiếu " + (player.combineNew.gemCombine - player.inventory.gem) + " ngọc xanh\n";
                    CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU, npcSay, "Đóng");
                } else if (player.inventory.gold < player.combineNew.goldCombine) {
                    npcSay += "|7|Còn thiếu " + Util.powerToString(player.combineNew.goldCombine - player.inventory.gold) + " vàng\n";
                    CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU, npcSay, "Đóng");
                } else {
                    CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.MENU_START_COMBINE, npcSay,
                            "Nâng cấp\n" + Util.numberToMoney(player.combineNew.goldCombine) + " vàng\n"
                            + Util.numberToMoney(player.combineNew.gemCombine) + " ngọc\n", "Từ chối");
                }
            } else {
                CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                        "Cần 1 Sao Pha Lê cấp 1 và 1 đá Hematite", "Đóng");
            }
        } else {
            CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                    "Cần 1 Sao Pha Lê cấp 1 và 1 đá Hematite", "Đóng");
        }
    }

    // Nâng cấp sao Pha Lê
    public static void nangCapSaoPhaLe(Player player) {
        if (player.combineNew.itemsCombine.size() == 2) {
            int gold = player.combineNew.goldCombine;
            int gem = player.combineNew.gemCombine;

            if (player.inventory.gold < gold) {
                Service.gI().sendThongBao(player, "Không đủ vàng để thực hiện");
                return;
            }
            if (player.inventory.gem < gem) {
                Service.gI().sendThongBao(player, "Không đủ ngọc để thực hiện");
                return;
            }

            Item saoPhaLe = null;
            Item hematite = null;

            // Tìm kiếm sao Pha Lê và Hematite
            for (Item item : player.combineNew.itemsCombine) {
                if (item.template.id >= 441 && item.template.id <= 447) {
                    saoPhaLe = item; // Sao Pha Lê cấp 1
                } else if (item.template.id == 1423) {
                    hematite = item; // Hematite
                }
            }

            if (saoPhaLe != null && hematite != null) {
                player.inventory.gold -= gold;
                player.inventory.gem -= gem;
                
                if (Util.isTrue(player.combineNew.ratioCombine, 100)) {
                    int getSaoPhaLeCap2Id = getSaoPhaLeCap2Id(saoPhaLe.template.id);
                    ItemTemplate newTemplate = ItemService.gI().getTemplate(getSaoPhaLeCap2Id);
                    Item newItem = new Item();
                    newItem.template = newTemplate;
                    newItem.quantity = 1;
                    newItem.itemOptions.clear();
                    
                    // Sao chép các options từ sao Pha Lê cấp 1 và cộng thêm 1 vào giá trị của option
                    for (Item.ItemOption option : saoPhaLe.itemOptions) {
                        newItem.itemOptions.add(new Item.ItemOption(option.optionTemplate.id, option.param));
                    }

                    // Thêm sao Pha Lê cấp 2 vào hành trang
                    InventoryService.gI().addItemBag(player, newItem);
                    
                    // Tiêu thụ sao Pha Lê cấp 1
                    InventoryService.gI().subQuantityItemsBag(player, saoPhaLe, 1);
                    CombineService.gI().sendEffectSuccessCombine(player);
                } else {
                    CombineService.gI().sendEffectFailCombine(player);
                }

                // Tiêu thụ đá Hematite
                InventoryService.gI().subQuantityItemsBag(player, hematite, 1); 
                InventoryService.gI().sendItemBags(player);
                Service.gI().sendMoney(player);
                CombineService.gI().reOpenItemCombine(player);
            }
        }
    }

    // Lấy ID của sao Pha Lê cấp 2 từ ID sao Pha Lê cấp 1
    private static int getSaoPhaLeCap2Id(int saoPhaLeCap1Id) {
        switch (saoPhaLeCap1Id) {
            case 441:
                return 1416; // Sao pha lê đỏ
            case 442:
                return 1417; // Sao pha lê lam
            case 443:
                return 1418; // Sao pha lê hồng
            case 444:
                return 1419; // Sao pha lê tím
            case 445:
                return 1420; // Sao pha lê cam
            case 446:
                return 1421; // Sao pha lê vàng
            case 447:
                return 1422; // Sao pha lê lục
            default:
                return -1;
        }
    }
}

